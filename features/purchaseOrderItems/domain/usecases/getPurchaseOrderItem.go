package purchaseorderitempresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getPurchaseOrderItemUsecase struct {
	getPurchaseOrderItemRepository purchaseorderitemdomainrepositoryinterfaces.GetPurchaseOrderItemRepository
}

func NewGetPurchaseOrderItemUsecase(
	getPurchaseOrderItemRepository purchaseorderitemdomainrepositoryinterfaces.GetPurchaseOrderItemRepository,
) (purchaseorderitempresentationusecaseinterfaces.GetPurchaseOrderItemUsecase, error) {
	return &getPurchaseOrderItemUsecase{
		getPurchaseOrderItemRepository,
	}, nil
}

func (getPurchaseOrderItemUcase *getPurchaseOrderItemUsecase) validation(
	input *model.PurchaseOrderItemFilterFields,
) (*model.PurchaseOrderItemFilterFields, error) {
	return input, nil
}

func (getPurchaseOrderItemUcase *getPurchaseOrderItemUsecase) Execute(
	filterFields *model.PurchaseOrderItemFilterFields,
) (*model.PurchaseOrderItem, error) {
	validatedFilterFields, err := getPurchaseOrderItemUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	purchaseOrderItem, err := getPurchaseOrderItemUcase.getPurchaseOrderItemRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPurchaseOrderItem",
			err,
		)
	}
	return purchaseOrderItem, nil
}
