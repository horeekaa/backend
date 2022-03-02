package paymentpresentationusecases

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
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updatePaymentUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatePaymentRepo    paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository
	approveUpdatePaymentRepo    paymentdomainrepositoryinterfaces.ApproveUpdatePaymentRepository
	updatePaymentAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdatePaymentUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdatePaymentRepo paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository,
	approveUpdatePaymentRepo paymentdomainrepositoryinterfaces.ApproveUpdatePaymentRepository,
) (paymentpresentationusecaseinterfaces.UpdatePaymentUsecase, error) {
	return &updatePaymentUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdatePaymentRepo,
		approveUpdatePaymentRepo,
		&model.MemberAccessRefOptionsInput{
			PaymentAccesses: &model.PaymentAccessesInput{
				PaymentUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updatePaymentUcase *updatePaymentUsecase) validation(input paymentpresentationusecasetypes.UpdatePaymentUsecaseInput) (paymentpresentationusecasetypes.UpdatePaymentUsecaseInput, error) {
	if &input.Context == nil {
		return paymentpresentationusecasetypes.UpdatePaymentUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updatePaymentUsecase",
				nil,
			)
	}

	return input, nil
}

func (updatePaymentUcase *updatePaymentUsecase) Execute(input paymentpresentationusecasetypes.UpdatePaymentUsecaseInput) (*model.Payment, error) {
	validatedInput, err := updatePaymentUcase.validation(input)
	if err != nil {
		return nil, err
	}
	paymentToUpdate := &model.InternalUpdatePayment{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdatePayment)
	json.Unmarshal(jsonTemp, paymentToUpdate)

	account, err := updatePaymentUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePaymentUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updatePaymentUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updatePaymentUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updatePaymentUcase.updatePaymentAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePaymentUsecase",
			err,
		)
	}

	jsonTemp, _ = json.Marshal(accMemberAccess)
	json.Unmarshal(jsonTemp, &paymentToUpdate.MemberAccess)

	if validatedInput.UpdatePayment.Photo != nil {
		paymentToUpdate.Photo.Photo.File = validatedInput.UpdatePayment.Photo.Photo.File
	}

	// if user is only going to approve proposal
	if paymentToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.PaymentAccesses.PaymentApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updatePaymentUsecase",
				nil,
			)
		}
		if !*accMemberAccess.Access.PaymentAccesses.PaymentApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updatePaymentUsecase",
				nil,
			)
		}

		paymentToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updatePaymentOutput, err := updatePaymentUcase.approveUpdatePaymentRepo.RunTransaction(
			paymentToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updatePaymentUsecase",
				err,
			)
		}

		return updatePaymentOutput, nil
	}

	paymentToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.PaymentAccesses.PaymentApproval != nil {
		if *accMemberAccess.Access.PaymentAccesses.PaymentApproval {
			paymentToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	paymentToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updatePaymentOutput, err := updatePaymentUcase.proposeUpdatePaymentRepo.RunTransaction(
		paymentToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePaymentUsecase",
			err,
		)
	}

	return updatePaymentOutput, nil
}
