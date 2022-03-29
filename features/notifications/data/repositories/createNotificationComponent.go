package notificationdomainrepositories

import (
	"context"
	"encoding/json"
	"time"

	"firebase.google.com/go/v4/messaging"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type createNotificationTransactionComponent struct {
	notificationDataSource   databasenotificationdatasourceinterfaces.NotificationDataSource
	firebaseMessaging        firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation
	notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder
	masterPayloadLoader      notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader
	pathIdentity             string
}

func NewCreateNotificationTransactionComponent(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
	notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
	masterPayloadLoader notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader,
) (notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent, error) {
	return &createNotificationTransactionComponent{
		notificationDataSource:   notificationDataSource,
		firebaseMessaging:        firebaseMessaging,
		notifLocalizationBuilder: notifLocalizationBuilder,
		masterPayloadLoader:      masterPayloadLoader,
		pathIdentity:             "CreateNotificationComponent",
	}, nil
}

func (createNotifTrx *createNotificationTransactionComponent) PreTransaction(
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

	_, err := createNotifTrx.masterPayloadLoader.TransactionBody(
		session,
		notificationToCreate,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createNotifTrx.pathIdentity,
			err,
		)
	}
	currentTime := time.Now()
	notificationToCreate.CreatedAt = &currentTime
	notificationToCreate.UpdatedAt = &currentTime

	createdNotification, err := createNotifTrx.notificationDataSource.GetMongoDataSource().Create(
		notificationToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createNotifTrx.pathIdentity,
			err,
		)
	}

	notificationToOutput := &model.Notification{}
	json.Unmarshal(jsonTemp, notificationToOutput)

	createNotifTrx.notifLocalizationBuilder.Execute(
		createdNotification,
		notificationToOutput,
		createdNotification.RecipientAccount.Language.String(),
	)

	payloadJson, _ := json.Marshal(notificationToCreate.PayloadOptions)
	_, err = createNotifTrx.firebaseMessaging.SendMulticastMessage(
		context.Background(),
		&firebasemessagingcoretypes.SentMulticastMessage{
			Notification: &messaging.Notification{
				Title: notificationToOutput.Message.Title,
				Body:  notificationToOutput.Message.Body,
			},
			Data: map[string]string{
				"payloadOptions": string(payloadJson),
			},
			Tokens: createdNotification.RecipientAccount.DeviceTokens,
		},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createNotifTrx.pathIdentity,
			err,
		)
	}

	return notificationToOutput, nil
}
