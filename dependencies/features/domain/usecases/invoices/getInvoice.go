package invoicepresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicepresentationusecases "github.com/horeekaa/backend/features/invoices/domain/usecases"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
)

type GetInvoiceUsecaseDependency struct{}

func (_ GetInvoiceUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getInvoiceRepo invoicedomainrepositoryinterfaces.GetInvoiceRepository,
		) invoicepresentationusecaseinterfaces.GetInvoiceUsecase {
			getInvoiceUsecase, _ := invoicepresentationusecases.NewGetInvoiceUsecase(
				getInvoiceRepo,
			)
			return getInvoiceUsecase
		},
	)
}
