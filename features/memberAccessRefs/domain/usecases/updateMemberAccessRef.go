package memberaccessrefpresentationusecases

import (
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessRefUsecase struct {
	getAccountFromAuthDataRepo          accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo          memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo        accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	updateMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository
	getMemberAccessRefRepo              memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository
	logEntityProposalActivityRepo       loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	logEntityApprovalActivityRepo       loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository
	updateMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	updateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository,
	getMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
) (memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase, error) {
	return &updateMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		updateMemberAccessRefRepo,
		getMemberAccessRefRepo,
		logEntityProposalActivityRepo,
		logEntityApprovalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			MemberAccessRefAccesses: &model.MemberAccessRefAccessesInput{
				MemberAccessRefUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateMmbAccessRefUcase *updateMemberAccessRefUsecase) validation(input memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput) (memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput, error) {
	if &input.Context == nil {
		return memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateMemberAccessRefUsecase",
				nil,
			)
	}
	input.UpdateMemberAccessRef.ApprovingAccount = nil
	input.UpdateMemberAccessRef.SubmittingAccount = nil

	return input, nil
}

func (updateMmbAccessRefUcase *updateMemberAccessRefUsecase) Execute(input memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput) (*model.MemberAccessRef, error) {
	validatedInput, err := updateMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateMmbAccessRefUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updateMemberAccessRefUsecase",
			nil,
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
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *updateMmbAccessRefUcase.updateMemberAccessRefAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
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
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	existingMemberAccRef, err := updateMmbAccessRefUcase.getMemberAccessRefRepo.Execute(
		&model.MemberAccessRefFilterFields{
			ID: &validatedInput.UpdateMemberAccessRef.ID,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	// if user is only going to approve proposal
	if validatedInput.UpdateMemberAccessRef.ProposalStatus != nil {
		if accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateMemberAccessRefUsecase",
				nil,
			)
		}
		if !*accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateMemberAccessRefUsecase",
				nil,
			)
		}

		logApprovalActivity, err := updateMmbAccessRefUcase.logEntityApprovalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityApprovalActivityInput{
				PreviousLog:      existingMemberAccRef.CorrespondingLog,
				ApprovingAccount: account,
				ApproverInitial:  accountInitials,
				ApprovalStatus:   *validatedInput.UpdateMemberAccessRef.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessRefUsecase",
				err,
			)
		}

		validatedInput.UpdateMemberAccessRef.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		validatedInput.UpdateMemberAccessRef.CorrespondingLog = &model.ObjectIDOnly{ID: &logApprovalActivity.ID}
		updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.updateMemberAccessRefRepo.RunTransaction(
			validatedInput.UpdateMemberAccessRef,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessRefUsecase",
				err,
			)
		}

		return updateMemberAccessRefOutput.UpdatedMemberAccessRef, nil
	}

	validatedInput.UpdateMemberAccessRef.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval != nil {
		if *accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval {
			validatedInput.UpdateMemberAccessRef.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	var newObject interface{} = *validatedInput.UpdateMemberAccessRef
	var existingObject interface{} = *existingMemberAccRef
	logEntityProposal, err := updateMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "MemberAccessRef",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityUpdate,
			ProposalStatus:   *validatedInput.UpdateMemberAccessRef.ProposalStatus,
			NewObject:        &newObject,
			ExistingObject:   &existingObject,
			ExistingObjectID: func(t string) *string { return &t }(existingMemberAccRef.ID.Hex()),
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	validatedInput.UpdateMemberAccessRef.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.UpdateMemberAccessRef.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.updateMemberAccessRefRepo.RunTransaction(
		validatedInput.UpdateMemberAccessRef,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	// user is going to update and directly has permission to approve the update
	if *validatedInput.UpdateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
		updateMemberAccessRefOutput, err = updateMmbAccessRefUcase.updateMemberAccessRefRepo.RunTransaction(
			&model.UpdateMemberAccessRef{
				ID:               updateMemberAccessRefOutput.UpdatedMemberAccessRef.ID,
				ApprovingAccount: &model.ObjectIDOnly{ID: &account.ID},
				ProposalStatus:   validatedInput.UpdateMemberAccessRef.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessRefUsecase",
				err,
			)
		}

		return updateMemberAccessRefOutput.UpdatedMemberAccessRef, nil
	}

	return updateMemberAccessRefOutput.UpdatedMemberAccessRef, nil
}
