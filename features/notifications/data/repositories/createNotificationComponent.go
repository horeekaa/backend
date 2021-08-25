package notificationdomainrepositories

import (
	"context"
	"encoding/json"

	"firebase.google.com/go/v4/messaging"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createNotificationTransactionComponent struct {
	notificationDataSource   databasenotificationdatasourceinterfaces.NotificationDataSource
	firebaseMessaging        firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder
}

func NewCreateNotificationTransactionComponent(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
) (notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent, error) {
	return &createNotificationTransactionComponent{
		notificationDataSource:   notificationDataSource,
		firebaseMessaging:        firebaseMessaging,
		notifLocalizationBuilder: notifLocalizationBuilder,
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

	notificationToOutput := &model.Notification{}
	json.Unmarshal(jsonTemp, notificationToOutput)

	createNotifTrx.notifLocalizationBuilder.Execute(createdNotification, notificationToOutput)

	_, err = createNotifTrx.firebaseMessaging.SendMulticastMessage(
		context.Background(),
		&firebasemessagingcoretypes.SentMulticastMessage{
			Notification: &messaging.Notification{
				Title: notificationToOutput.Message.Title,
				Body:  notificationToOutput.Message.Body,
			},
			Tokens: input.RecipientAccount.DeviceTokens,
		},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createNotification",
			err,
		)
	}

	return notificationToOutput, nil
}
