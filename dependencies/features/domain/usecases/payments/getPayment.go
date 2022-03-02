package paymentpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentpresentationusecases "github.com/horeekaa/backend/features/payments/domain/usecases"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
)

type GetPaymentUsecaseDependency struct{}

func (_ GetPaymentUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getPaymentRepo paymentdomainrepositoryinterfaces.GetPaymentRepository,
		) paymentpresentationusecaseinterfaces.GetPaymentUsecase {
			getPaymentUsecase, _ := paymentpresentationusecases.NewGetPaymentUsecase(
				getPaymentRepo,
			)
			return getPaymentUsecase
		},
	)
}
