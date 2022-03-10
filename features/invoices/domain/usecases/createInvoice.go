package invoicepresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type createInvoiceUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createInvoiceRepo           invoicedomainrepositoryinterfaces.CreateInvoiceRepository
	createInvoiceAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewCreateInvoiceUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createinvoiceRepo invoicedomainrepositoryinterfaces.CreateInvoiceRepository,
) (invoicepresentationusecaseinterfaces.CreateInvoiceUsecase, error) {
	return &createInvoiceUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createinvoiceRepo,
		&model.MemberAccessRefOptionsInput{
			InvoiceAccesses: &model.InvoiceAccessesInput{
				InvoiceCreate: func(b bool) *bool { return &b }(true),
			},
		},
		"CreateInvoiceUsecase",
	}, nil
}

func (createInvoiceUcase *createInvoiceUsecase) validation(input invoicepresentationusecasetypes.CreateInvoiceUsecaseInput) (invoicepresentationusecasetypes.CreateInvoiceUsecaseInput, error) {
	if &input.Context == nil {
		return invoicepresentationusecasetypes.CreateInvoiceUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				createInvoiceUcase.pathIdentity,
				nil,
			)
	}
	return input, nil
}

func (createInvoiceUcase *createInvoiceUsecase) Execute(input invoicepresentationusecasetypes.CreateInvoiceUsecaseInput) ([]*model.Invoice, error) {
	validatedInput, err := createInvoiceUcase.validation(input)
	if err != nil {
		return nil, err
	}

	if !validatedInput.CronAuthenticated {
		account, err := createInvoiceUcase.getAccountFromAuthDataRepo.Execute(
			accountdomainrepositorytypes.GetAccountFromAuthDataInput{
				Context: validatedInput.Context,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				createInvoiceUcase.pathIdentity,
				err,
			)
		}
		if account == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				createInvoiceUcase.pathIdentity,
				nil,
			)
		}

		memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
		_, err = createInvoiceUcase.getAccountMemberAccessRepo.Execute(
			memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
				MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
					Account:             &model.ObjectIDOnly{ID: &account.ID},
					MemberAccessRefType: &memberAccessRefTypeOrgBased,
					Access:              createInvoiceUcase.createInvoiceAccessIdentity,
				},
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				createInvoiceUcase.pathIdentity,
				err,
			)
		}
	}

	invoiceToCreate := &model.InternalCreateInvoice{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateInvoice)
	json.Unmarshal(jsonTemp, invoiceToCreate)

	createdInvoice, err := createInvoiceUcase.createInvoiceRepo.RunTransaction(
		invoiceToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			createInvoiceUcase.pathIdentity,
			err,
		)
	}

	return createdInvoice, nil
}
