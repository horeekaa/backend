package notificationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationpresentationusecases "github.com/horeekaa/backend/features/notifications/domain/usecases"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
)

type GetAllNotificationUsecaseDependency struct{}

func (_ GetAllNotificationUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllNotificationRepo notificationdomainrepositoryinterfaces.GetAllNotificationRepository,
		) notificationpresentationusecaseinterfaces.GetAllNotificationUsecase {
			getAllNotificationUcase, _ := notificationpresentationusecases.NewGetAllNotificationUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllNotificationRepo,
			)
			return getAllNotificationUcase
		},
	)
}
