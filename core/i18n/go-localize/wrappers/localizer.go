package golocalizei18ncorewrappers

import (
	golocalizei18ncorewrapperinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/wrappers"
	golocalizei18ncoretypes "github.com/horeekaa/backend/core/i18n/go-localize/types"
	localizations "github.com/horeekaa/backend/features/i18n/data/dataSources/generated"
)

type goLocalizeI18NLocalizer struct {
	*localizations.Localizer
}

func NewGoLocalizeI18NLocalizer(localizer *localizations.Localizer) (golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer, error) {
	return &goLocalizeI18NLocalizer{
		localizer,
	}, nil
}

func (goLocalizeLocalizer *goLocalizeI18NLocalizer) Get(key string, replacements ...*golocalizei18ncoretypes.LocalizerReplacement) string {
	nativeReplacements := []*localizations.Replacements{}

	for _, replacement := range replacements {
		nativeReplacement := replacement.Replacements
		nativeReplacements = append(nativeReplacements, nativeReplacement)
	}
	return goLocalizeLocalizer.Localizer.Get(key, nativeReplacements...)
}

func (goLocalizeLocalizer *goLocalizeI18NLocalizer) GetWithLocale(locale, key string, replacements ...*golocalizei18ncoretypes.LocalizerReplacement) string {
	nativeReplacements := []*localizations.Replacements{}

	for _, replacement := range replacements {
		nativeReplacement := replacement.Replacements
		nativeReplacements = append(nativeReplacements, nativeReplacement)
	}
	return goLocalizeLocalizer.Localizer.GetWithLocale(locale, key, nativeReplacements...)
}

func (goLocalizeLocalizer *goLocalizeI18NLocalizer) SetFallbackLocale(fallback string) golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer {
	goLocalizeLocalizer.Localizer.SetFallbackLocale(fallback)
	return goLocalizeLocalizer
}

func (goLocalizeLocalizer *goLocalizeI18NLocalizer) SetLocale(locale string) golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer {
	goLocalizeLocalizer.Localizer.SetLocale(locale)
	return goLocalizeLocalizer
}

func (goLocalizeLocalizer *goLocalizeI18NLocalizer) SetLocales(locale, fallback string) golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer {
	goLocalizeLocalizer.Localizer.SetLocales(locale, fallback)
	return goLocalizeLocalizer
}
