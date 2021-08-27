package notificationdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	notificationdomainrepositoryutilities "github.com/horeekaa/backend/features/notifications/data/repositories/utils"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
)

type NotificationLocalizationBuilderDependency struct{}

func (_ *NotificationLocalizationBuilderDependency) Bind() {
	container.Singleton(
		func(
			goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient,
		) notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder {
			notificationLocalizationBuilder, _ := notificationdomainrepositoryutilities.NewNotificationLocalizationBuilder(
				goLocalizeI18N,
			)
			return notificationLocalizationBuilder
		},
	)
}
