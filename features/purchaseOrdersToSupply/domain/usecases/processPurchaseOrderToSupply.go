package purchaseordertosupplypresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

type processPurchaseOrderToSupplyUsecase struct {
	processPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyRepository
	pathIdentity                     string
}

func NewProcessPurchaseOrderToSupplyUsecase(
	processPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyRepository,
) (purchaseordertosupplypresentationusecaseinterfaces.ProcessPurchaseOrderToSupplyUsecase, error) {
	return &processPurchaseOrderToSupplyUsecase{
		processPurchaseOrderToSupplyRepo: processPurchaseOrderToSupplyRepo,
		pathIdentity:                     "ProcessPurchaseOrderToSupplyUsecase",
	}, nil
}

func (processPurchaseOrderToSupplyUcase *processPurchaseOrderToSupplyUsecase) Execute() (bool, error) {
	ok, err := processPurchaseOrderToSupplyUcase.processPurchaseOrderToSupplyRepo.RunTransaction()
	if err != nil {
		return false, horeekaacorefailuretoerror.ConvertFailure(
			processPurchaseOrderToSupplyUcase.pathIdentity,
			err,
		)
	}

	return ok, nil
}
