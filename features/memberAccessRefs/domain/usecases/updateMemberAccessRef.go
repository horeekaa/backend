package memberaccessrefpresentationusecases

import (
	"encoding/json"

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
	updateMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository
	getMemberAccessRefRepo              memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository
	logEntityProposalActivityRepo       loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	logEntityApprovalActivityRepo       loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository
	updateMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	updateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository,
	getMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
) (memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase, error) {
	return &updateMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
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

	accMemberAccess, err := updateMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account: &model.ObjectIDOnly{ID: &account.ID},
				Access:  updateMmbAccessRefUcase.updateMemberAccessRefAccessIdentity,
			},
		},
	)
	if err != nil {
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

	memberAccessRefToUpdate := &model.InternalUpdateMemberAccessRef{
		ID: validatedInput.UpdateMemberAccessRef.ID,
	}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateMemberAccessRef)
	json.Unmarshal(jsonTemp, memberAccessRefToUpdate)

	// if user is only going to approve proposal
	if memberAccessRefToUpdate.ProposalStatus != nil {
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
				PreviousLog:      existingMemberAccRef.RecentLog,
				ApprovingAccount: account,
				ApprovalStatus:   *memberAccessRefToUpdate.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessRefUsecase",
				err,
			)
		}

		memberAccessRefToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		memberAccessRefToUpdate.RecentLog = &model.ObjectIDOnly{ID: &logApprovalActivity.ID}
		updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.updateMemberAccessRefRepo.RunTransaction(
			memberAccessRefToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessRefUsecase",
				err,
			)
		}

		return updateMemberAccessRefOutput, nil
	}

	memberAccessRefToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval != nil {
		if *accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval {
			memberAccessRefToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	var newObject interface{} = *memberAccessRefToUpdate
	var existingObject interface{} = *existingMemberAccRef
	logEntityProposal, err := updateMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "MemberAccessRef",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityUpdate,
			ProposalStatus:   *memberAccessRefToUpdate.ProposalStatus,
			NewObject:        &newObject,
			ExistingObject:   &existingObject,
			ExistingObjectID: &model.ObjectIDOnly{ID: &existingMemberAccRef.ID},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	memberAccessRefToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	memberAccessRefToUpdate.RecentLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.updateMemberAccessRefRepo.RunTransaction(
		memberAccessRefToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessRefUsecase",
			err,
		)
	}

	return updateMemberAccessRefOutput, nil
}
