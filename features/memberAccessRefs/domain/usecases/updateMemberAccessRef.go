package memberaccessrefpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
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
	proposeUpdateMemberAccessRefRepo    memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefRepository
	approveUpdateMemberAccessRefRepo    memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefRepository
	updateMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefRepository,
	approveUpdateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefRepository,
) (memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase, error) {
	return &updateMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateMemberAccessRefRepo,
		approveUpdateMemberAccessRefRepo,
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
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
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

		memberAccessRefToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.approveUpdateMemberAccessRefRepo.RunTransaction(
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

	memberAccessRefToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateMemberAccessRefOutput, err := updateMmbAccessRefUcase.proposeUpdateMemberAccessRefRepo.RunTransaction(
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
