package organizationpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createOrganizationUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createOrganizationRepo           organizationdomainrepositoryinterfaces.CreateOrganizationRepository
	createMemberAccessRepo           memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository
	createOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createOrganizationRepo organizationdomainrepositoryinterfaces.CreateOrganizationRepository,
	createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
) (organizationpresentationusecaseinterfaces.CreateOrganizationUsecase, error) {
	return &createOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createOrganizationRepo,
		createMemberAccessRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createOrganizationUcase *createOrganizationUsecase) validation(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (organizationpresentationusecasetypes.CreateOrganizationUsecaseInput, error) {
	if &input.Context == nil {
		return organizationpresentationusecasetypes.CreateOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createOrganizationUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateOrganization.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createOrganizationUcase *createOrganizationUsecase) Execute(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := createOrganizationUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createOrganizationUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createOrganizationUsecase",
			nil,
		)
	}

	memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
	accMemberAccess, err := createOrganizationUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
				Access:              createOrganizationUcase.createOrganizationAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}
	if accMemberAccess.Access.OrganizationAccesses.OrganizationApproval != nil {
		if *accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
			validatedInput.CreateOrganization.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	organizationToCreate := &model.InternalCreateOrganization{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateOrganization)
	json.Unmarshal(jsonTemp, organizationToCreate)

	organizationToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdOrganization, err := createOrganizationUcase.createOrganizationRepo.RunTransaction(
		organizationToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}
	if createdOrganization.ProposalStatus == model.EntityProposalStatusApproved {
		createOrganizationUcase.createMemberAccessRepo.RunTransaction(
			&model.InternalCreateMemberAccess{
				Account: &model.ObjectIDOnly{
					ID: &account.ID,
				},
				Organization: &model.InternalUpdateOrganization{
					ID: createdOrganization.ID,
				},
				OrganizationMembershipRole: func(r model.OrganizationMembershipRole) *model.OrganizationMembershipRole {
					return &r
				}(model.OrganizationMembershipRoleOwner),
				MemberAccessRefType: model.MemberAccessRefTypeOrganizationsBased,
				SubmittingAccount: &model.ObjectIDOnly{
					ID: &createdOrganization.RecentApprovingAccount.ID,
				},
				ProposalStatus: func(e model.EntityProposalStatus) *model.EntityProposalStatus {
					return &e
				}(model.EntityProposalStatusApproved),
				InvitationAccepted: func(b bool) *bool {
					return &b
				}(true),
			},
		)
	}

	return createdOrganization, nil
}
