package purchaseorderpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getPurchaseOrderUsecase struct {
	getPurchaseOrderRepository purchaseorderdomainrepositoryinterfaces.GetPurchaseOrderRepository
	pathIdentity               string
}

func NewGetPurchaseOrderUsecase(
	getPurchaseOrderRepository purchaseorderdomainrepositoryinterfaces.GetPurchaseOrderRepository,
) (purchaseorderpresentationusecaseinterfaces.GetPurchaseOrderUsecase, error) {
	return &getPurchaseOrderUsecase{
		getPurchaseOrderRepository,
		"GetPurchaseOrderUsecase",
	}, nil
}

func (getPOUcase *getPurchaseOrderUsecase) validation(
	input *model.PurchaseOrderFilterFields,
) (*model.PurchaseOrderFilterFields, error) {
	return input, nil
}

func (getPOUcase *getPurchaseOrderUsecase) Execute(
	filterFields *model.PurchaseOrderFilterFields,
) (*model.PurchaseOrder, error) {
	validatedFilterFields, err := getPOUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	purchaseOrder, err := getPOUcase.getPurchaseOrderRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getPOUcase.pathIdentity,
			err,
		)
	}
	return purchaseOrder, nil
}
