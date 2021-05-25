package memberaccesspresentationusecases

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
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessUsecase struct {
	getAccountFromAuthDataRepo    accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo    memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo  accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	updateMemberAccessRepo        memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository
	updateMemberAccessIdentity    *model.MemberAccessRefOptionsInput
}

func NewUpdateMemberAccessUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	updateMemberAccessRepo memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
) (memberaccesspresentationusecaseinterfaces.UpdateMemberAccessUsecase, error) {
	return &updateMemberAccessUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		updateMemberAccessRepo,
		logEntityProposalActivityRepo,
		logEntityApprovalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateMmbAccessUcase *updateMemberAccessUsecase) validation(input memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput) (memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput, error) {
	if &input.Context == nil {
		return memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateMemberAccessUsecase",
				nil,
			)
	}
	input.UpdateMemberAccess.ApprovingAccount = nil
	input.UpdateMemberAccess.SubmittingAccount = nil

	return input, nil
}

func (updateMmbAccessUcase *updateMemberAccessUsecase) Execute(input memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput) (*model.MemberAccess, error) {
	validatedInput, err := updateMmbAccessUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateMmbAccessUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updateMemberAccessUsecase",
			nil,
		)
	}

	personChannel := make(chan *model.Person)
	errChannel := make(chan error)
	go func() {
		person, err := updateMmbAccessUcase.getPersonDataFromAccountRepo.Execute(account)
		if err != nil {
			errChannel <- err
		}
		personChannel <- person
	}()

	accMemberAccess, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *updateMmbAccessUcase.updateMemberAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
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
			"/updateMemberAccessUsecase",
			err,
		)
	}

	existingMemberAcc, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			Account: &model.Account{}, //validatedInput.UpdateMemberAccess.ID,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}

	// if user is only going to approve proposal
	if validatedInput.UpdateMemberAccess.ProposalStatus != nil {
		if accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateMemberAccessUsecase",
				nil,
			)
		}
		if !*accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateMemberAccessUsecase",
				nil,
			)
		}

		logApprovalActivity, err := updateMmbAccessUcase.logEntityApprovalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityApprovalActivityInput{
				PreviousLog:      existingMemberAcc.CorrespondingLog,
				ApprovingAccount: account,
				ApproverInitial:  accountInitials,
				ApprovalStatus:   *validatedInput.UpdateMemberAccess.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
				err,
			)
		}

		validatedInput.UpdateMemberAccess.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		validatedInput.UpdateMemberAccess.CorrespondingLog = &model.ObjectIDOnly{ID: &logApprovalActivity.ID}
		updateMemberAccessOutput, err := updateMmbAccessUcase.updateMemberAccessRepo.RunTransaction(
			validatedInput.UpdateMemberAccess,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
				err,
			)
		}

		return updateMemberAccessOutput.UpdatedMemberAccess, nil
	}

	validatedInput.UpdateMemberAccess.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval != nil {
		if *accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval {
			validatedInput.UpdateMemberAccess.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	var newObject interface{} = *validatedInput.UpdateMemberAccess
	var existingObject interface{} = *existingMemberAcc
	logEntityProposal, err := updateMmbAccessUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "MemberAccess",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityUpdate,
			ProposalStatus:   *validatedInput.UpdateMemberAccess.ProposalStatus,
			NewObject:        &newObject,
			ExistingObject:   &existingObject,
			ExistingObjectID: func(t string) *string { return &t }(existingMemberAcc.ID.Hex()),
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}

	validatedInput.UpdateMemberAccess.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.UpdateMemberAccess.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	updateMemberAccessOutput, err := updateMmbAccessUcase.updateMemberAccessRepo.RunTransaction(
		validatedInput.UpdateMemberAccess,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}

	// user is going to update and directly has permission to approve the update
	if *validatedInput.UpdateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
		updateMemberAccessOutput, err = updateMmbAccessUcase.updateMemberAccessRepo.RunTransaction(
			&model.UpdateMemberAccess{
				ID:               updateMemberAccessOutput.UpdatedMemberAccess.ID,
				ApprovingAccount: &model.ObjectIDOnly{ID: &account.ID},
				ProposalStatus:   validatedInput.UpdateMemberAccess.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
				err,
			)
		}

		return updateMemberAccessOutput.UpdatedMemberAccess, nil
	}

	return updateMemberAccessOutput.UpdatedMemberAccess, nil
}
