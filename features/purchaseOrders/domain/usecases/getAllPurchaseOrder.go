package purchaseorderpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/types"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllPurchaseOrderUsecase struct {
	getAccountFromAuthDataRepo        accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo        memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllpurchaseOrderRepo           purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository
	getAllpurchaseOrderAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllPurchaseOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllpurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository,
) (purchaseorderpresentationusecaseinterfaces.GetAllPurchaseOrderUsecase, error) {
	return &getAllPurchaseOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllpurchaseOrderRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllPOUcase *getAllPurchaseOrderUsecase) validation(input purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput) (*purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput, error) {
	if &input.Context == nil {
		return &purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllPurchaseOrderUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllPOUcase *getAllPurchaseOrderUsecase) Execute(
	input purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput,
) ([]*model.PurchaseOrder, error) {
	validatedInput, err := getAllPOUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllPOUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllPurchaseOrderUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllPOUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllPOUcase.getAllpurchaseOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}

	purchaseOrders, err := getAllPOUcase.getAllpurchaseOrderRepo.Execute(
		purchaseorderdomainrepositorytypes.GetAllPurchaseOrderInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}

	return purchaseOrders, nil
}
