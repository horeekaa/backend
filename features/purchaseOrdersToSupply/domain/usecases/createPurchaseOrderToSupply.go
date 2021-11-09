package purchaseordertosupplypresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderToSupplyUsecase struct {
	createPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyRepository
}

func NewCreatePurchaseOrderToSupplyUsecase(
	createpurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyRepository,
) (purchaseordertosupplypresentationusecaseinterfaces.CreatePurchaseOrderToSupplyUsecase, error) {
	return &createPurchaseOrderToSupplyUsecase{
		createPurchaseOrderToSupplyRepo: createpurchaseOrderToSupplyRepo,
	}, nil
}

func (createPurchaseOrderToSupplyUcase *createPurchaseOrderToSupplyUsecase) Execute() ([]*model.PurchaseOrderToSupply, error) {
	purchaseOrdersToSupply, err := createPurchaseOrderToSupplyUcase.createPurchaseOrderToSupplyRepo.RunTransaction()
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPurchaseOrderToSupplyUsecase",
			err,
		)
	}

	return purchaseOrdersToSupply, nil
}
