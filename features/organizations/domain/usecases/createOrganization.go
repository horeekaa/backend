package organizationpresentationusecases

import (
	"errors"
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createOrganizationUsecase struct {
	manageAccountAuthenticationRepo  accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository
	getAccountMemberAccessRepo       accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	createOrganizationRepo           organizationdomainrepositoryinterfaces.CreateOrganizationRepository
	logEntityProposalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	createOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateOrganizationUsecase(
	manageAccountAuthenticationRepo accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	createOrganizationRepo organizationdomainrepositoryinterfaces.CreateOrganizationRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
) (organizationpresentationusecaseinterfaces.CreateOrganizationUsecase, error) {
	return &createOrganizationUsecase{
		manageAccountAuthenticationRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		createOrganizationRepo,
		logEntityProposalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMmbAccessRefUcase *createOrganizationUsecase) validation(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (organizationpresentationusecasetypes.CreateOrganizationUsecaseInput, error) {
	if &input.AuthHeader == nil {
		return organizationpresentationusecasetypes.CreateOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createOrganizationUsecase",
				errors.New(horeekaacoreerrorenums.AuthenticationTokenNotExist),
			)
	}
	input.CreateOrganization.SubmittingAccount = nil
	input.CreateOrganization.ProposalStatus = nil
	return input, nil
}

func (createMmbAccessRefUcase *createOrganizationUsecase) Execute(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := createMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMmbAccessRefUcase.manageAccountAuthenticationRepo.RunTransaction(
		accountdomainrepositorytypes.ManageAccountAuthenticationInput{
			AuthHeader: validatedInput.AuthHeader,
			Context:    validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
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

	accMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		accountdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *createMmbAccessRefUcase.createOrganizationAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createOrganizationUsecase",
			err,
		)
	}
	if *accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
		validatedInput.CreateOrganization.ProposalStatus =
			func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
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

	return createdOrganization, nil
}
