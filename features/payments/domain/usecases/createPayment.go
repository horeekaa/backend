package paymentpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createPaymentUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createPaymentRepo           paymentdomainrepositoryinterfaces.CreatePaymentRepository
	createPaymentAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreatePaymentUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createPaymentRepo paymentdomainrepositoryinterfaces.CreatePaymentRepository,
) (paymentpresentationusecaseinterfaces.CreatePaymentUsecase, error) {
	return &createPaymentUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createPaymentRepo,
		&model.MemberAccessRefOptionsInput{
			PaymentAccesses: &model.PaymentAccessesInput{
				PaymentCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createPaymentUcase *createPaymentUsecase) validation(input paymentpresentationusecasetypes.CreatePaymentUsecaseInput) (paymentpresentationusecasetypes.CreatePaymentUsecaseInput, error) {
	if &input.Context == nil {
		return paymentpresentationusecasetypes.CreatePaymentUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createPaymentUsecase",
				nil,
			)
	}
	return input, nil
}

func (createPaymentUcase *createPaymentUsecase) Execute(input paymentpresentationusecasetypes.CreatePaymentUsecaseInput) (*model.Payment, error) {
	validatedInput, err := createPaymentUcase.validation(input)
	if err != nil {
		return nil, err
	}
	paymentToCreate := &model.InternalCreatePayment{}
	jsonTemp, _ := json.Marshal(validatedInput.CreatePayment)
	json.Unmarshal(jsonTemp, paymentToCreate)

	account, err := createPaymentUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPaymentUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createPaymentUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createPaymentUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createPaymentUcase.createPaymentAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPaymentUsecase",
			err,
		)
	}

	paymentToCreate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.PaymentAccesses.PaymentApproval != nil {
		if *accMemberAccess.Access.PaymentAccesses.PaymentApproval {
			paymentToCreate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	jsonTemp, _ = json.Marshal(accMemberAccess)
	json.Unmarshal(jsonTemp, &paymentToCreate.MemberAccess)

	if validatedInput.CreatePayment.Photo != nil {
		paymentToCreate.Photo.Photo.File = validatedInput.CreatePayment.Photo.Photo.File
	}

	paymentToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdPayment, err := createPaymentUcase.createPaymentRepo.RunTransaction(
		paymentToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPaymentUsecase",
			err,
		)
	}

	return createdPayment, nil
}
