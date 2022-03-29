package notificationpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositorytypes "github.com/horeekaa/backend/features/notifications/domain/repositories/types"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllNotificationUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllNotificationRepo           notificationdomainrepositoryinterfaces.GetAllNotificationRepository
	getAllNotificationAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                     string
}

func NewGetAllNotificationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllNotificationRepo notificationdomainrepositoryinterfaces.GetAllNotificationRepository,
) (notificationpresentationusecaseinterfaces.GetAllNotificationUsecase, error) {
	return &getAllNotificationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllNotificationRepo,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountReadOwned: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllNotificationUsecase",
	}, nil
}

func (getAllNotificationUcase *getAllNotificationUsecase) validation(input notificationpresentationusecasetypes.GetAllNotificationUsecaseInput) (*notificationpresentationusecasetypes.GetAllNotificationUsecaseInput, error) {
	if &input.Context == nil {
		return &notificationpresentationusecasetypes.GetAllNotificationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllNotificationUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllNotificationUcase *getAllNotificationUsecase) Execute(
	input notificationpresentationusecasetypes.GetAllNotificationUsecaseInput,
) ([]*model.Notification, error) {
	validatedInput, err := getAllNotificationUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllNotificationUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllNotificationUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllNotificationUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
	_, err = getAllNotificationUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
				Access:              getAllNotificationUcase.getAllNotificationAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllNotificationUcase.pathIdentity,
			err,
		)
	}

	notifications, err := getAllNotificationUcase.getAllNotificationRepo.Execute(
		notificationdomainrepositorytypes.GetAllNotificationInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
			Language:      account.Language.String(),
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllNotificationUcase.pathIdentity,
			err,
		)
	}

	return notifications, nil
}
