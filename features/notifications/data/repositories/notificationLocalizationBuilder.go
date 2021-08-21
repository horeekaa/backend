package notificationdomainrepositories

import (
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	golocalizei18ncoretypes "github.com/horeekaa/backend/core/i18n/go-localize/types"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type notificationLocalizationBuilder struct {
	goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient
}

func NewNotificationLocalizationBuilder(
	goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient,
) (notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder, error) {
	return &notificationLocalizationBuilder{
		goLocalizeI18N: goLocalizeI18N,
	}, nil
}

func (notifLocalBuilder *notificationLocalizationBuilder) Execute(
	input *model.InternalCreateNotification,
	output *model.Notification,
) (bool, error) {
	notifLocalBuilder.goLocalizeI18N.Initialize(
		input.RecipientAccount.Language.String(),
		"id",
	)
	localizer, _ := notifLocalBuilder.goLocalizeI18N.GetLocalizer()

	titleText := ""
	bodyText := ""
	switch input.NotificationCategory {
	case model.NotificationCategoryOrgInvitationAccepted:
		titleText = localizer.Get(
			"organizations.invitationAccepted.messages.invitation_accepted_notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"personName": input.PayloadOptions.InvitationAcceptedPayload.MemberAccess.Account.Person.FirstName,
			},
		)
		bodyText = localizer.Get(
			"organizations.invitationAccepted.messages.invitation_accepted_notification_body",
		)
		break

	case model.NotificationCategoryOrgInvitationRequest:
		titleText = localizer.Get(
			"organizations.invitationRequest.messages.invitation_request_notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"submitterName": input.PayloadOptions.InvitationRequestPayload.MemberAccess.SubmittingAccount.Person.FirstName,
				"orgName":       input.PayloadOptions.InvitationRequestPayload.MemberAccess.Organization.Name,
			},
		)
		bodyText = localizer.Get(
			"organizations.invitationRequest.messages.invitation_request_notification_body",
		)
		break
	}

	(*output).Message = &model.NotificationMessage{
		Title: titleText,
		Body:  bodyText,
	}

	return true, nil
}
