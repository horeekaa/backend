package notificationdomainrepositories

import (
	"context"
	"encoding/json"

	"firebase.google.com/go/v4/messaging"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type createNotificationTransactionComponent struct {
	accountDataSource        databaseaccountdatasourceinterfaces.AccountDataSource
	notificationDataSource   databasenotificationdatasourceinterfaces.NotificationDataSource
	firebaseMessaging        firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation
	notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder
	masterPayloadLoader      notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader
}

func NewCreateNotificationTransactionComponent(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
	notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
	masterPayloadLoader notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader,
) (notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent, error) {
	return &createNotificationTransactionComponent{
		accountDataSource:        accountDataSource,
		notificationDataSource:   notificationDataSource,
		firebaseMessaging:        firebaseMessaging,
		notifLocalizationBuilder: notifLocalizationBuilder,
		masterPayloadLoader:      masterPayloadLoader,
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
	_, err := createNotifTrx.masterPayloadLoader.TransactionBody(
		session,
		input,
	)
	if err != nil {
		return nil, err
	}

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

	recipientAccount, err := createNotifTrx.accountDataSource.GetMongoDataSource().FindByID(
		*input.RecipientAccount.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createNotification",
			err,
		)
	}

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
			Tokens: recipientAccount.DeviceTokens,
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
