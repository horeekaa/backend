package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type NotificationLocalizationBuilderDependency struct{}

func (_ *NotificationLocalizationBuilderDependency) Bind() {
	container.Singleton(
		func(
			goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient,
		) notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder {
			notificationLocalizationBuilder, _ := notificationdomainrepositories.NewNotificationLocalizationBuilder(
				goLocalizeI18N,
			)
			return notificationLocalizationBuilder
		},
	)
}
