package golocalizei18ncoreclientinterfaces

import (
	golocalizei18ncorewrapperinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/wrappers"
)

type GoLocalizeI18NClient interface {
	Initialize(locale string, fallbackLocale string) (bool, error)
	GetLocalizer() (golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer, error)
}
