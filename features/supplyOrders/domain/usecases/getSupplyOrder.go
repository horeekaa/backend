package supplyorderpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getSupplyOrderUsecase struct {
	getSupplyOrderRepository supplyorderdomainrepositoryinterfaces.GetSupplyOrderRepository
	pathIdentity             string
}

func NewGetSupplyOrderUsecase(
	getSupplyOrderRepository supplyorderdomainrepositoryinterfaces.GetSupplyOrderRepository,
) (supplyorderpresentationusecaseinterfaces.GetSupplyOrderUsecase, error) {
	return &getSupplyOrderUsecase{
		getSupplyOrderRepository,
		"GetSupplyOrderUsecase",
	}, nil
}

func (getSOUcase *getSupplyOrderUsecase) validation(
	input *model.SupplyOrderFilterFields,
) (*model.SupplyOrderFilterFields, error) {
	return input, nil
}

func (getSOUcase *getSupplyOrderUsecase) Execute(
	filterFields *model.SupplyOrderFilterFields,
) (*model.SupplyOrder, error) {
	validatedFilterFields, err := getSOUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	supplyOrder, err := getSOUcase.getSupplyOrderRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getSOUcase.pathIdentity,
			err,
		)
	}
	return supplyOrder, nil
}
