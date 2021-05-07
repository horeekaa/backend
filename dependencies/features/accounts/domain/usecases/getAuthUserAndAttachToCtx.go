package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type GetAuthUserAndAttachToCtxUsecaseDependency struct{}

func (_ GetAuthUserAndAttachToCtxUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getUserFromAuthHeaderRepo accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository,
		) accountpresentationusecaseinterfaces.GetAuthUserAndAttachToCtxUsecase {
			getUsrAndAttachToCtxUcase, _ := accountpresentationusecases.NewGetAuthUserAndAttachToCtxUsecase(
				getUserFromAuthHeaderRepo,
			)

			return getUsrAndAttachToCtxUcase
		},
	)
}
