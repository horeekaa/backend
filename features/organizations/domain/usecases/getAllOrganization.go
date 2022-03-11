package organizationpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllOrganizationUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllOrganizationRepo     organizationdomainrepositoryinterfaces.GetAllOrganizationRepository

	getOwnedOrganizationIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                 string
}

func NewGetAllOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllOrganizationRepo organizationdomainrepositoryinterfaces.GetAllOrganizationRepository,
) (organizationpresentationusecaseinterfaces.GetAllOrganizationUsecase, error) {
	return &getAllOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllOrganizationRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationReadOwned: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllOrganizationUsecase",
	}, nil
}

func (getAllOrgUcase *getAllOrganizationUsecase) validation(input organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput) (*organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput, error) {
	if &input.Context == nil {
		return &organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllOrgUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllOrgUcase *getAllOrganizationUsecase) Execute(
	input organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput,
) ([]*model.Organization, error) {
	validatedInput, err := getAllOrgUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllOrgUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllOrgUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllOrgUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllOrgUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Status: func(s model.MemberAccessStatus) *model.MemberAccessStatus {
					return &s
				}(model.MemberAccessStatusActive),
				ProposalStatus: func(e model.EntityProposalStatus) *model.EntityProposalStatus {
					return &e
				}(model.EntityProposalStatusApproved),
				InvitationAccepted: func(b bool) *bool {
					return &b
				}(true),
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllOrgUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.OrganizationAccesses.OrganizationReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.OrganizationAccesses.OrganizationReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.ID = &memberAccess.Organization.ID
		} else {
			memberAccessRefTypeAccountBasics := model.MemberAccessRefTypeAccountsBasics
			_, err := getAllOrgUcase.getAccountMemberAccessRepo.Execute(
				memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
					MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
						Account:             &model.ObjectIDOnly{ID: &account.ID},
						MemberAccessRefType: &memberAccessRefTypeAccountBasics,
						Access:              getAllOrgUcase.getOwnedOrganizationIdentity,
					},
				},
			)
			if err != nil {
				return nil, horeekaacorefailuretoerror.ConvertFailure(
					getAllOrgUcase.pathIdentity,
					err,
				)
			}
			validatedInput.FilterFields.SubmittingAccount = &model.ObjectIDOnly{
				ID: &account.ID,
			}
		}
	}

	organizations, err := getAllOrgUcase.getAllOrganizationRepo.Execute(
		organizationdomainrepositorytypes.GetAllOrganizationInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllOrgUcase.pathIdentity,
			err,
		)
	}

	return organizations, nil
}
