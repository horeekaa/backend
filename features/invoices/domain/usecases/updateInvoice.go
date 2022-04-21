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

type updateInvoiceUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	updateInvoiceRepo           invoicedomainrepositoryinterfaces.UpdateInvoiceRepository
	updateDueInvoiceRepo        invoicedomainrepositoryinterfaces.UpdateDueInvoiceRepository
	updateInvoiceAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewUpdateInvoiceUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	updateInvoiceRepo invoicedomainrepositoryinterfaces.UpdateInvoiceRepository,
	updateDueInvoiceRepo invoicedomainrepositoryinterfaces.UpdateDueInvoiceRepository,
) (invoicepresentationusecaseinterfaces.UpdateInvoiceUsecase, error) {
	return &updateInvoiceUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		updateInvoiceRepo,
		updateDueInvoiceRepo,
		&model.MemberAccessRefOptionsInput{
			InvoiceAccesses: &model.InvoiceAccessesInput{
				InvoiceUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateInvoiceUsecase",
	}, nil
}

func (updateInvoiceUcase *updateInvoiceUsecase) validation(input invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput) (invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput, error) {
	if &input.Context == nil {
		return invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateInvoiceUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateInvoiceUcase *updateInvoiceUsecase) Execute(input invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput) (*model.Invoice, error) {
	validatedInput, err := updateInvoiceUcase.validation(input)
	if err != nil {
		return nil, err
	}

	if validatedInput.CronAuthenticated {
		_, err = updateInvoiceUcase.updateDueInvoiceRepo.RunTransaction()
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateInvoiceUcase.pathIdentity,
				err,
			)
		}
		return nil, nil
	}

	account, err := updateInvoiceUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateInvoiceUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateInvoiceUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = updateInvoiceUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateInvoiceUcase.updateInvoiceAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateInvoiceUcase.pathIdentity,
			err,
		)
	}

	invoiceToUpdate := &model.InternalUpdateInvoice{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateInvoice)
	json.Unmarshal(jsonTemp, invoiceToUpdate)

	updateInvoiceOutput, err := updateInvoiceUcase.updateInvoiceRepo.RunTransaction(
		invoiceToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateInvoiceUcase.pathIdentity,
			err,
		)
	}

	return updateInvoiceOutput, nil
}
