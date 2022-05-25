package notificationpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createNotificationMessage struct {
	notifLocalizationBuilder   notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	pathIdentity               string
}

func NewCreateNotificationMessage(
	notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
) (notificationpresentationusecaseinterfaces.CreateNotificationMessage, error) {
	return &createNotificationMessage{
		notifLocalizationBuilder,
		getAccountFromAuthDataRepo,
		"CreateNotificationMessage",
	}, nil
}

func (createNotifMessageUcase *createNotificationMessage) validation(input notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput) (notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput, error) {
	if &input.Context == nil {
		return notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				createNotifMessageUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (createNotifMessageUcase *createNotificationMessage) Execute(input notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput) (*model.NotificationMessage, error) {
	validatedInput, err := createNotifMessageUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createNotifMessageUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			createNotifMessageUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			createNotifMessageUcase.pathIdentity,
			nil,
		)
	}

	databaseNotification := &model.DatabaseNotification{}
	jsonOutput, _ := json.Marshal(validatedInput.Notification)
	json.Unmarshal(jsonOutput, databaseNotification)

	createNotifMessageUcase.notifLocalizationBuilder.Execute(
		databaseNotification,
		validatedInput.Notification,
		account.Language.String(),
	)

	return validatedInput.Notification.Message, nil
}
