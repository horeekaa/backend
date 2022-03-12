package purchaseordertosupplypresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplydomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories/types"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	purchaseordertosupplypresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllPurchaseOrderToSupplyUsecase struct {
	getAccountFromAuthDataRepo        accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo        memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllPurchaseOrderToSupplyRepo   purchaseordertosupplydomainrepositoryinterfaces.GetAllPurchaseOrderToSupplyRepository
	getAllPurchaseOrderAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                      string
}

func NewGetAllPurchaseOrderToSupplyUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.GetAllPurchaseOrderToSupplyRepository,
) (purchaseordertosupplypresentationusecaseinterfaces.GetAllPurchaseOrderToSupplyUsecase, error) {
	return &getAllPurchaseOrderToSupplyUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllPurchaseOrderToSupplyRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllPurchaseOrderToSupplyUsecase",
	}, nil
}

func (getAllPOToSupplyUcase *getAllPurchaseOrderToSupplyUsecase) validation(input purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput) (*purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput, error) {
	if &input.Context == nil {
		return &purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllPOToSupplyUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllPOToSupplyUcase *getAllPurchaseOrderToSupplyUsecase) Execute(
	input purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput,
) ([]*model.PurchaseOrderToSupply, error) {
	validatedInput, err := getAllPOToSupplyUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllPOToSupplyUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllPOToSupplyUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllPOToSupplyUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllPOToSupplyUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllPOToSupplyUcase.getAllPurchaseOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllPOToSupplyUcase.pathIdentity,
			err,
		)
	}

	purchaseOrdersToSupply, err := getAllPOToSupplyUcase.getAllPurchaseOrderToSupplyRepo.Execute(
		purchaseordertosupplydomainrepositorytypes.GetAllPurchaseOrderToSupplyInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllPOToSupplyUcase.pathIdentity,
			err,
		)
	}

	return purchaseOrdersToSupply, nil
}
