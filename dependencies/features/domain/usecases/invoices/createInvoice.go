package invoicepresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicepresentationusecases "github.com/horeekaa/backend/features/invoices/domain/usecases"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type CreateInvoiceUsecaseDependency struct{}

func (_ *CreateInvoiceUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createInvoiceRepo invoicedomainrepositoryinterfaces.CreateInvoiceRepository,
		) invoicepresentationusecaseinterfaces.CreateInvoiceUsecase {
			invoiceRefUcase, _ := invoicepresentationusecases.NewCreateInvoiceUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createInvoiceRepo,
			)
			return invoiceRefUcase
		},
	)
}
