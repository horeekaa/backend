package organizationpresentationusecases

import (
	"errors"
	"fmt"

	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/serviceFailures/enums"
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/usecaseErrors"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/usecaseErrors/_failureToError"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/usecaseErrors/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateOrganizationUsecase struct {
	manageAccountAuthenticationRepo  accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository
	getAccountMemberAccessRepo       accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	updateOrganizationRepo           organizationdomainrepositoryinterfaces.UpdateOrganizationRepository
	getOrganizationRepo              organizationdomainrepositoryinterfaces.GetOrganizationRepository
	logEntityProposalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	logEntityApprovalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository
	updateOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateOrganizationUsecase(
	manageAccountAuthenticationRepo accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	updateOrganizationRepo organizationdomainrepositoryinterfaces.UpdateOrganizationRepository,
	getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
) (organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase, error) {
	return &updateOrganizationUsecase{
		manageAccountAuthenticationRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		updateOrganizationRepo,
		getOrganizationRepo,
		logEntityProposalActivityRepo,
		logEntityApprovalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateMmbAccessRefUcase *updateOrganizationUsecase) validation(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput, error) {
	if &input.AuthHeader == nil {
		return organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateOrganizationUsecase",
				errors.New(horeekaacoreerrorenums.AuthenticationTokenNotExist),
			)
	}
	input.UpdateOrganization.ApprovingAccount = nil
	input.UpdateOrganization.SubmittingAccount = nil

	return input, nil
}

func (updateMmbAccessRefUcase *updateOrganizationUsecase) Execute(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := updateMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateMmbAccessRefUcase.manageAccountAuthenticationRepo.RunTransaction(
		accountdomainrepositorytypes.ManageAccountAuthenticationInput{
			AuthHeader: validatedInput.AuthHeader,
			Context:    validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	personChannel := make(chan *model.Person)
	errChannel := make(chan error)
	go func() {
		person, err := updateMmbAccessRefUcase.getPersonDataFromAccountRepo.Execute(account)
		if err != nil {
			errChannel <- err
		}
		personChannel <- person
	}()

	accMemberAccess, err := updateMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		accountdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *updateMmbAccessRefUcase.updateOrganizationAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
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
			"/updateOrganizationUsecase",
			err,
		)
	}

	existingMemberAccRef, err := updateMmbAccessRefUcase.getOrganizationRepo.Execute(
		&model.OrganizationFilterFields{
			ID: &validatedInput.UpdateOrganization.ID,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	// if user is only going to approve proposal
	if validatedInput.UpdateOrganization.ProposalStatus != nil {
		if !*accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateOrganizationUsecase",
				errors.New(horeekaacorefailureenums.FeatureNotAccessibleByAccount),
			)
		}

		logApprovalActivity, err := updateMmbAccessRefUcase.logEntityApprovalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityApprovalActivityInput{
				PreviousLog:      existingMemberAccRef.CorrespondingLog,
				ApprovingAccount: account,
				ApproverInitial:  accountInitials,
				ApprovalStatus:   *validatedInput.UpdateOrganization.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		validatedInput.UpdateOrganization.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		validatedInput.UpdateOrganization.CorrespondingLog = &model.ObjectIDOnly{ID: &logApprovalActivity.ID}
		updateOrganizationOutput, err := updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
			validatedInput.UpdateOrganization,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		return updateOrganizationOutput.UpdatedOrganization, nil
	}

	if *accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
		validatedInput.UpdateOrganization.ProposalStatus =
			func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
	}

	var newObject interface{} = *validatedInput.UpdateOrganization
	var existingObject interface{} = *existingMemberAccRef
	logEntityProposal, err := updateMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "Organization",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityUpdate,
			ProposalStatus:   *validatedInput.UpdateOrganization.ProposalStatus,
			NewObject:        &newObject,
			ExistingObject:   &existingObject,
			ExistingObjectID: func(t string) *string { return &t }(existingMemberAccRef.ID.Hex()),
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	validatedInput.UpdateOrganization.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.UpdateOrganization.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	updateOrganizationOutput, err := updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
		validatedInput.UpdateOrganization,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	// user is going to update and directly has permission to approve the update
	if *validatedInput.UpdateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
		updateOrganizationOutput, err = updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
			&model.UpdateOrganization{
				ID:               updateOrganizationOutput.UpdatedOrganization.ID,
				ApprovingAccount: &model.ObjectIDOnly{ID: &account.ID},
				ProposalStatus:   validatedInput.UpdateOrganization.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		return updateOrganizationOutput.UpdatedOrganization, nil
	}

	return updateOrganizationOutput.UpdatedOrganization, nil
}
