package notificationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationpresentationusecases "github.com/horeekaa/backend/features/notifications/domain/usecases"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
)

type CreateNotificationMessageUsecaseDependency struct{}

func (_ *CreateNotificationMessageUsecaseDependency) Bind() {
	container.Singleton(
		func(
			notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
		) notificationpresentationusecaseinterfaces.CreateNotificationMessage {
			createNotificationMessageUsecase, _ := notificationpresentationusecases.NewCreateNotificationMessage(
				notifLocalizationBuilder,
				getAccountFromAuthDataRepo,
			)
			return createNotificationMessageUsecase
		},
	)
}
