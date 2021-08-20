package i18ncoredependencies

import (
	"github.com/golobby/container/v2"
	golocalizei18ncoreclients "github.com/horeekaa/backend/core/i18n/go-localize"
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
)

type GoLocalizeI18NDependency struct{}

func (_ GoLocalizeI18NDependency) Bind() {
	container.Singleton(
		func() golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient {
			goLocalizeClient, _ := golocalizei18ncoreclients.NewGoLocalizeI18NClient()
			return goLocalizeClient
		},
	)
}
