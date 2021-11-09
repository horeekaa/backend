package purchaseordertosupplypresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getPurchaseOrderToSupplyUsecase struct {
	getPurchaseOrderToSupplyRepository purchaseordertosupplydomainrepositoryinterfaces.GetPurchaseOrderToSupplyRepository
}

func NewGetPurchaseOrderToSupplyUsecase(
	getPurchaseOrderToSupplyRepository purchaseordertosupplydomainrepositoryinterfaces.GetPurchaseOrderToSupplyRepository,
) (purchaseordertosupplypresentationusecaseinterfaces.GetPurchaseOrderToSupplyUsecase, error) {
	return &getPurchaseOrderToSupplyUsecase{
		getPurchaseOrderToSupplyRepository,
	}, nil
}

func (getPOToSupplyUcase *getPurchaseOrderToSupplyUsecase) validation(
	input *model.PurchaseOrderToSupplyFilterFields,
) (*model.PurchaseOrderToSupplyFilterFields, error) {
	return input, nil
}

func (getPOToSupplyUcase *getPurchaseOrderToSupplyUsecase) Execute(
	filterFields *model.PurchaseOrderToSupplyFilterFields,
) (*model.PurchaseOrderToSupply, error) {
	validatedFilterFields, err := getPOToSupplyUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	purchaseOrderToSupply, err := getPOToSupplyUcase.getPurchaseOrderToSupplyRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPurchaseOrderToSupply",
			err,
		)
	}
	return purchaseOrderToSupply, nil
}
