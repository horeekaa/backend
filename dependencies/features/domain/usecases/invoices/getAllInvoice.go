package invoicepresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicepresentationusecases "github.com/horeekaa/backend/features/invoices/domain/usecases"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAllInvoiceUsecaseDependency struct{}

func (_ GetAllInvoiceUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountinvoicepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllInvoiceRepo invoicedomainrepositoryinterfaces.GetAllInvoiceRepository,
		) invoicepresentationusecaseinterfaces.GetAllInvoiceUsecase {
			getAllInvoiceUcase, _ := invoicepresentationusecases.NewGetAllInvoiceUsecase(
				getAccountFromAuthDataRepo,
				getAccountinvoicepo,
				getAllInvoiceRepo,
			)
			return getAllInvoiceUcase
		},
	)
}
