package invoicepresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getInvoiceUsecase struct {
	getInvoiceRepository invoicedomainrepositoryinterfaces.GetInvoiceRepository
}

func NewGetInvoiceUsecase(
	getInvoiceRepository invoicedomainrepositoryinterfaces.GetInvoiceRepository,
) (invoicepresentationusecaseinterfaces.GetInvoiceUsecase, error) {
	return &getInvoiceUsecase{
		getInvoiceRepository,
	}, nil
}

func (getInvoiceUcase *getInvoiceUsecase) validation(
	input *model.InvoiceFilterFields,
) (*model.InvoiceFilterFields, error) {
	return input, nil
}

func (getInvoiceUcase *getInvoiceUsecase) Execute(
	filterFields *model.InvoiceFilterFields,
) (*model.Invoice, error) {
	validatedFilterFields, err := getInvoiceUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	invoice, err := getInvoiceUcase.getInvoiceRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getInvoiceUsecase",
			err,
		)
	}
	return invoice, nil
}
