package memberaccesspresentationusecases

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
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessUsecase struct {
	getAccountFromAuthDataRepo           accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo           memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateMemberAccessRepo        memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository
	approveUpdateMemberAccessRepo        memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository
	updateMemberAccessIdentity           *model.MemberAccessRefOptionsInput
	acceptInvitationMemberAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                         string
}

func NewUpdateMemberAccessUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateMemberAccessRepo memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository,
	approveUpdateMemberAccessRepo memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository,
) (memberaccesspresentationusecaseinterfaces.UpdateMemberAccessUsecase, error) {
	return &updateMemberAccessUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateMemberAccessRepo,
		approveUpdateMemberAccessRepo,
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessAcceptInvitation: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateMemberAccessUsecase",
	}, nil
}

func (updateMmbAccessUcase *updateMemberAccessUsecase) validation(input memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput) (memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput, error) {
	if &input.Context == nil {
		return memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateMmbAccessUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateMmbAccessUcase *updateMemberAccessUsecase) Execute(input memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput) (*model.MemberAccess, error) {
	validatedInput, err := updateMmbAccessUcase.validation(input)
	if err != nil {
		return nil, err
	}
	memberAccessToUpdate := &model.InternalUpdateMemberAccess{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateMemberAccess)
	json.Unmarshal(jsonTemp, memberAccessToUpdate)

	account, err := updateMmbAccessUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMmbAccessUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateMmbAccessUcase.pathIdentity,
			nil,
		)
	}

	existingMemberAcc, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				ID: &memberAccessToUpdate.ID,
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMmbAccessUcase.pathIdentity,
			err,
		)
	}

	if memberAccessToUpdate.InvitationAccepted != nil {
		memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
		_, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
			memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
				MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
					Account:             &model.ObjectIDOnly{ID: &account.ID},
					MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
					Access:              updateMmbAccessUcase.acceptInvitationMemberAccessIdentity,
				},
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateMmbAccessUcase.pathIdentity,
				err,
			)
		}

		memberAccessToUpdate.SubmittingAccount = &model.ObjectIDOnly{
			ID: &account.ID,
		}
		memberAccessToUpdate.ProposalStatus = &existingMemberAcc.ProposalStatus
		updateMemberAccessOutput, err := updateMmbAccessUcase.proposeUpdateMemberAccessRepo.RunTransaction(
			memberAccessToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateMmbAccessUcase.pathIdentity,
				err,
			)
		}
		return updateMemberAccessOutput, nil
	}

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Access:              updateMmbAccessUcase.updateMemberAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMmbAccessUcase.pathIdentity,
			err,
		)
	}

	// if user is only going to approve proposal
	if memberAccessToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateMmbAccessUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateMmbAccessUcase.pathIdentity,
				nil,
			)
		}

		memberAccessToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateMemberAccessOutput, err := updateMmbAccessUcase.approveUpdateMemberAccessRepo.RunTransaction(
			memberAccessToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateMmbAccessUcase.pathIdentity,
				err,
			)
		}

		return updateMemberAccessOutput, nil
	}

	memberAccessToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval != nil {
		if *accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval {
			memberAccessToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	memberAccessToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateMemberAccessOutput, err := updateMmbAccessUcase.proposeUpdateMemberAccessRepo.RunTransaction(
		memberAccessToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateMmbAccessUcase.pathIdentity,
			err,
		)
	}

	return updateMemberAccessOutput, nil
}
