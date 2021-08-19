package golocalizei18ncorewrapperinterfaces

import (
	golocalizei18ncoretypes "github.com/horeekaa/backend/core/i18n/go-localize/types"
)

type GoLocalizeI18NLocalizer interface {
	Get(key string, replacements ...*golocalizei18ncoretypes.LocalizerReplacement) string
	GetWithLocale(locale, key string, replacements ...*golocalizei18ncoretypes.LocalizerReplacement) string
	SetFallbackLocale(fallback string) GoLocalizeI18NLocalizer
	SetLocale(locale string) GoLocalizeI18NLocalizer
	SetLocales(locale, fallback string) GoLocalizeI18NLocalizer
}
