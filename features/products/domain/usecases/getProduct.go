package productpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getProductUsecase struct {
	getProductRepository productdomainrepositoryinterfaces.GetProductRepository
}

func NewGetProductUsecase(
	getProductRepository productdomainrepositoryinterfaces.GetProductRepository,
) (productpresentationusecaseinterfaces.GetProductUsecase, error) {
	return &getProductUsecase{
		getProductRepository,
	}, nil
}

func (getProdUcase *getProductUsecase) validation(
	input *model.ProductFilterFields,
) (*model.ProductFilterFields, error) {
	return input, nil
}

func (getProdUcase *getProductUsecase) Execute(
	filterFields *model.ProductFilterFields,
) (*model.Product, error) {
	validatedFilterFields, err := getProdUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	product, err := getProdUcase.getProductRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getProduct",
			err,
		)
	}
	return product, nil
}
