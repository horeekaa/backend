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

	existingMemberAcc, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				ID: &validatedInput.UpdateMemberAccess.ID,
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}

	memberAccessToUpdate := &model.InternalUpdateMemberAccess{
		ID: validatedInput.UpdateMemberAccess.ID,
	}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateMemberAccess)
	json.Unmarshal(jsonTemp, memberAccessToUpdate)

	if memberAccessToUpdate.InvitationAccepted != nil {
		if existingMemberAcc.Account.ID.Hex() != account.ID.Hex() ||
			existingMemberAcc.ProposalStatus != model.EntityProposalStatusApproved {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AcceptInvitationNotAllowed,
				422,
				"/updateMemberAccessUsecase",
				nil,
			)
		}
		memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
		_, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
			memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
				MemberAccessFilterFields: &model.MemberAccessFilterFields{
					Account:             &model.ObjectIDOnly{ID: &account.ID},
					MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
					Access:              updateMmbAccessUcase.acceptInvitationMemberAccessIdentity,
				},
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
				err,
			)
		}

		memberAccessToUpdate.SubmittingAccount = &model.ObjectIDOnly{
			ID: &existingMemberAcc.SubmittingAccount.ID,
		}
		memberAccessToUpdate.ProposalStatus = &existingMemberAcc.ProposalStatus
		updateMemberAccessOutput, err := updateMmbAccessUcase.proposeUpdateMemberAccessRepo.RunTransaction(
			memberAccessToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
				err,
			)
		}
		return updateMemberAccessOutput, nil
	}

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Access:              updateMmbAccessUcase.updateMemberAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateMemberAccessUsecase",
			err,
		)
	}

	// if user is only going to approve proposal
	if memberAccessToUpdate.ProposalStatus != nil {
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

		memberAccessToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateMemberAccessOutput, err := updateMmbAccessUcase.approveUpdateMemberAccessRepo.RunTransaction(
			memberAccessToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateMemberAccessUsecase",
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
			"/updateMemberAccessUsecase",
			err,
		)
	}

	return updateMemberAccessOutput, nil
}
