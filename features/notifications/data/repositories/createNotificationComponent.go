package notificationdomainrepositories

import (
	"context"
	"encoding/json"

	"firebase.google.com/go/v4/messaging"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	golocalizei18ncoretypes "github.com/horeekaa/backend/core/i18n/go-localize/types"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createNotificationTransactionComponent struct {
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource
	firebaseMessaging      firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation
	goLocalizeI18N         golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient
}

func NewCreateNotificationTransactionComponent(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
	goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient,
) (notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent, error) {
	return &createNotificationTransactionComponent{
		notificationDataSource: notificationDataSource,
		firebaseMessaging:      firebaseMessaging,
		goLocalizeI18N:         goLocalizeI18N,
	}, nil
}

func (_ *createNotificationTransactionComponent) PreTransaction(
	input *model.InternalCreateNotification,
) (*model.InternalCreateNotification, error) {
	return input, nil
}

func (createNotifTrx *createNotificationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateNotification,
) (*model.Notification, error) {
	notificationToCreate := &model.DatabaseCreateNotification{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, notificationToCreate)

	createdNotification, err := createNotifTrx.notificationDataSource.GetMongoDataSource().Create(
		notificationToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createNotification",
			err,
		)
	}

	createNotifTrx.goLocalizeI18N.Initialize(
		input.RecipientAccount.Language.String(),
		"id",
	)
	localizer, _ := createNotifTrx.goLocalizeI18N.GetLocalizer()

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

	createNotifTrx.firebaseMessaging.SendMulticastMessage(
		context.Background(),
		&firebasemessagingcoretypes.SentMulticastMessage{
			Notification: &messaging.Notification{
				Title: titleText,
				Body:  bodyText,
			},
			Tokens: input.RecipientAccount.DeviceTokens,
		},
	)

	notificationToOutput := &model.Notification{}
	jsonTemp, _ = json.Marshal(createdNotification)
	json.Unmarshal(jsonTemp, notificationToOutput)
	notificationToOutput.Message = &model.NotificationMessage{
		Title: titleText,
		Body:  bodyText,
	}

	return notificationToOutput, nil
}
