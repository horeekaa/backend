package golocalizei18ncoreclients

import (
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	golocalizei18ncorewrapperinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/wrappers"
	golocalizei18ncorewrappers "github.com/horeekaa/backend/core/i18n/go-localize/wrappers"
	localizations "github.com/horeekaa/backend/features/i18n/data/dataSources/generated"
)

type goLocalizeI18NClient struct {
	localizer golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer
}

func NewGoLocalizeI18NClient() (golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient, error) {
	return &goLocalizeI18NClient{}, nil
}

func (goLocalizeClient *goLocalizeI18NClient) Initialize(locale string, fallbackLocale string) (bool, error) {
	wrappedLocalizer, _ := golocalizei18ncorewrappers.NewGoLocalizeI18NLocalizer(
		localizations.New(locale, fallbackLocale),
	)

	goLocalizeClient.localizer = wrappedLocalizer
	return true, nil
}

func (goLocalizeI18NClient *goLocalizeI18NClient) GetLocalizer() (golocalizei18ncorewrapperinterfaces.GoLocalizeI18NLocalizer, error) {
	if goLocalizeI18NClient.localizer == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newGoLocalizeI18N",
			nil,
		)
	}
	return goLocalizeI18NClient.localizer, nil
}
