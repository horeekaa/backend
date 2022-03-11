package moupresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateMouUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateMouRepo       moudomainrepositoryinterfaces.ProposeUpdateMouRepository
	approveUpdateMouRepo       moudomainrepositoryinterfaces.ApproveUpdateMouRepository
	updateMouAccessIdentity    *model.MemberAccessRefOptionsInput
	pathIdentity               string
}

func NewUpdateMouUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateMouRepo moudomainrepositoryinterfaces.ProposeUpdateMouRepository,
	approveUpdateMouRepo moudomainrepositoryinterfaces.ApproveUpdateMouRepository,
) (moupresentationusecaseinterfaces.UpdateMouUsecase, error) {
	return &updateMouUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateMouRepo,
		approveUpdateMouRepo,
		&model.MemberAccessRefOptionsInput{
			MouAccesses: &model.MouAccessesInput{
				MouUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateMouUsecase",
	}, nil
}

func (updateMouUcase *updateMouUsecase) validation(input moupresentationusecasetypes.UpdateMouUsecaseInput) (moupresentationusecasetypes.UpdateMouUsecaseInput, error) {
	if &input.Context == nil {
		return moupresentationusecasetypes.UpdateMouUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateMouUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateMouUcase *updateMouUsecase) Execute(input moupresentationusecasetypes.UpdateMouUsecaseInput) (*model.Mou, error) {
	validatedInput, err := updateMouUcase.validation(input)
	if err != nil {
		return nil, err
	}
	mouToUpdate := &model.InternalUpdateMou{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateMou)
	json.Unmarshal(jsonTemp, mouToUpdate)

	account, err := updateMouUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMouUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateMouUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateMouUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateMouUcase.updateMouAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMouUcase.pathIdentity,
			err,
		)
	}

	// if user is only going to approve proposal
	if mouToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.MouAccesses.MouApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateMouUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.MouAccesses.MouApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateMouUcase.pathIdentity,
				nil,
			)
		}

		mouToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateMouOutput, err := updateMouUcase.approveUpdateMouRepo.RunTransaction(
			mouToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateMouUcase.pathIdentity,
				err,
			)
		}

		return updateMouOutput, nil
	}

	mouToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.MouAccesses.MouApproval != nil {
		if *accMemberAccess.Access.MouAccesses.MouApproval {
			mouToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	mouToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateMouOutput, err := updateMouUcase.proposeUpdateMouRepo.RunTransaction(
		mouToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMouUcase.pathIdentity,
			err,
		)
	}

	return updateMouOutput, nil
}
