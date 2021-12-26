package supplyorderitempresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getSupplyOrderItemUsecase struct {
	getSupplyOrderItemRepository supplyorderitemdomainrepositoryinterfaces.GetSupplyOrderItemRepository
}

func NewGetSupplyOrderItemUsecase(
	getSupplyOrderItemRepository supplyorderitemdomainrepositoryinterfaces.GetSupplyOrderItemRepository,
) (supplyorderitempresentationusecaseinterfaces.GetSupplyOrderItemUsecase, error) {
	return &getSupplyOrderItemUsecase{
		getSupplyOrderItemRepository,
	}, nil
}

func (getSupplyOrderItemUcase *getSupplyOrderItemUsecase) validation(
	input *model.SupplyOrderItemFilterFields,
) (*model.SupplyOrderItemFilterFields, error) {
	return input, nil
}

func (getSupplyOrderItemUcase *getSupplyOrderItemUsecase) Execute(
	filterFields *model.SupplyOrderItemFilterFields,
) (*model.SupplyOrderItem, error) {
	validatedFilterFields, err := getSupplyOrderItemUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	supplyOrderItem, err := getSupplyOrderItemUcase.getSupplyOrderItemRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getSupplyOrderItem",
			err,
		)
	}
	return supplyOrderItem, nil
}
