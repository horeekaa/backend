package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type GetAccountUsecaseDependency struct{}

func (_ GetAccountUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountRepo accountdomainrepositoryinterfaces.GetAccountRepository,
		) accountpresentationusecaseinterfaces.GetAccountUsecase {
			getAccountUsecase, _ := accountpresentationusecases.NewGetAccountUsecase(
				getAccountRepo,
			)
			return getAccountUsecase
		},
	)
}
