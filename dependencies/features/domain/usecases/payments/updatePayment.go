package paymentpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentpresentationusecases "github.com/horeekaa/backend/features/payments/domain/usecases"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
)

type UpdatePaymentUsecaseDependency struct{}

func (_ *UpdatePaymentUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdatePaymentRepo paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository,
			approveUpdatePaymentRepo paymentdomainrepositoryinterfaces.ApproveUpdatePaymentRepository,
		) paymentpresentationusecaseinterfaces.UpdatePaymentUsecase {
			updatePaymentUsecase, _ := paymentpresentationusecases.NewUpdatePaymentUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdatePaymentRepo,
				approveUpdatePaymentRepo,
			)
			return updatePaymentUsecase
		},
	)
}
