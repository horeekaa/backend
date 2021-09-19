package notificationpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type bulkUpdateNotificationUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	bulkUpdateNotificationRepo       notificationdomainrepositoryinterfaces.BulkUpdateNotificationRepository
	updateNotificationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewBulkUpdateNotificationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	bulkUpdateNotificationRepo notificationdomainrepositoryinterfaces.BulkUpdateNotificationRepository,
) (notificationpresentationusecaseinterfaces.BulkUpdateNotificationUsecase, error) {
	return &bulkUpdateNotificationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		bulkUpdateNotificationRepo,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateNotificationUcase *bulkUpdateNotificationUsecase) validation(input notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput) (notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput, error) {
	if &input.Context == nil {
		return notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/bulkUpdateNotificationUsecase",
				nil,
			)
	}

	return input, nil
}

func (updateNotificationUcase *bulkUpdateNotificationUsecase) Execute(input notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput) ([]*model.Notification, error) {
	validatedInput, err := updateNotificationUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateNotificationUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/bulkUpdateNotificationUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/bulkUpdateNotificationUsecase",
			nil,
		)
	}

	memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
	_, err = updateNotificationUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
				Access:              updateNotificationUcase.updateNotificationAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/bulkUpdateNotificationUsecase",
			err,
		)
	}

	notificationToUpdate := &model.InternalBulkUpdateNotification{
		IDs: validatedInput.BulkUpdateNotification.IDs,
	}
	jsonTemp, _ := json.Marshal(validatedInput.BulkUpdateNotification)
	json.Unmarshal(jsonTemp, notificationToUpdate)

	updatedNotifications, err := updateNotificationUcase.bulkUpdateNotificationRepo.RunTransaction(
		notificationToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/bulkUpdateNotificationUsecase",
			err,
		)
	}

	return updatedNotifications, nil
}
