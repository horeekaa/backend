package organizationpresentationusecases

import (
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
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
	getPersonDataFromAccountRepo     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	createOrganizationRepo           organizationdomainrepositoryinterfaces.CreateOrganizationRepository
	createMemberAccessRepo           memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository
	logEntityProposalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	createOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	createOrganizationRepo organizationdomainrepositoryinterfaces.CreateOrganizationRepository,
	createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
) (organizationpresentationusecaseinterfaces.CreateOrganizationUsecase, error) {
	return &createOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		createOrganizationRepo,
		createMemberAccessRepo,
		logEntityProposalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMmbAccessRefUcase *createOrganizationUsecase) validation(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (organizationpresentationusecasetypes.CreateOrganizationUsecaseInput, error) {
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
	input.CreateOrganization.SubmittingAccount = nil
	input.CreateOrganization.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createMmbAccessRefUcase *createOrganizationUsecase) Execute(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := createMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMmbAccessRefUcase.getAccountFromAuthDataRepo.Execute(
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

	personChannel := make(chan *model.Person)
	errChannel := make(chan error)
	go func() {
		person, err := createMmbAccessRefUcase.getPersonDataFromAccountRepo.Execute(account)
		if err != nil {
			errChannel <- err
		}
		personChannel <- person
	}()

	memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
	accMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
				Access:              createMmbAccessRefUcase.createOrganizationAccessIdentity,
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

	accountInitials := ""
	select {
	case person := <-personChannel:
		accountInitials = fmt.Sprintf("XXXX%s", account.ID.Hex()[len(account.ID.Hex())-6:])
		if person != nil {
			accountInitials = person.FirstName
		}

		break
	case err := <-errChannel:
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}

	var newObject interface{} = *validatedInput.CreateOrganization
	logEntityProposal, err := createMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "Organization",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityCreate,
			ProposalStatus:   *validatedInput.CreateOrganization.ProposalStatus,
			NewObject:        &newObject,
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}

	validatedInput.CreateOrganization.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.CreateOrganization.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	if *validatedInput.CreateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
		validatedInput.CreateOrganization.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
	}
	createdOrganization, err := createMmbAccessRefUcase.createOrganizationRepo.Execute(
		validatedInput.CreateOrganization,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}

	if createdOrganization.ProposalStatus == model.EntityProposalStatusApproved {
		newMemberAccess := &model.CreateMemberAccess{
			Account:            &model.ObjectIDOnly{ID: &account.ID},
			InvitationAccepted: func(b bool) *bool { return &b }(true),
			Organization: &model.AttachOrganizationInput{
				ID:   &createdOrganization.ID,
				Type: &createdOrganization.Type,
			},
			OrganizationMembershipRole: func(
				s model.OrganizationMembershipRole,
			) *model.OrganizationMembershipRole {
				return &s
			}(model.OrganizationMembershipRoleOwner),
			MemberAccessRefType: model.MemberAccessRefTypeOrganizationsBased,
			SubmittingAccount:   &model.ObjectIDOnly{ID: &account.ID},
			ApprovingAccount:    &model.ObjectIDOnly{ID: &account.ID},
			ProposalStatus: func(ep model.EntityProposalStatus) *model.EntityProposalStatus {
				return &ep
			}(model.EntityProposalStatusApproved),
		}

		var newObject interface{} = *newMemberAccess
		logEntityProposal, err := createMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityProposalActivityInput{
				CollectionName:   "MemberAccess",
				CreatedByAccount: account,
				Activity:         model.LoggedActivityCreate,
				ProposalStatus:   *newMemberAccess.ProposalStatus,
				NewObject:        &newObject,
				CreatorInitial:   accountInitials,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/createOrganizationUsecase",
				err,
			)
		}
		newMemberAccess.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}

		_, err = createMmbAccessRefUcase.createMemberAccessRepo.Execute(
			newMemberAccess,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/createOrganizationUsecase",
				err,
			)
		}
	}

	return createdOrganization, nil
}
