package productvariantpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	productvariantpresentationusecaseinterfaces "github.com/horeekaa/backend/features/productVariants/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getProductVariantUsecase struct {
	getProductVariantRepository productvariantdomainrepositoryinterfaces.GetProductVariantRepository
	pathIdentity                string
}

func NewGetProductVariantUsecase(
	getProductVariantRepository productvariantdomainrepositoryinterfaces.GetProductVariantRepository,
) (productvariantpresentationusecaseinterfaces.GetProductVariantUsecase, error) {
	return &getProductVariantUsecase{
		getProductVariantRepository,
		"GetProductVariantUsecase",
	}, nil
}

func (getProdVariantUcase *getProductVariantUsecase) validation(
	input *model.ProductVariantFilterFields,
) (*model.ProductVariantFilterFields, error) {
	return input, nil
}

func (getProdVariantUcase *getProductVariantUsecase) Execute(
	filterFields *model.ProductVariantFilterFields,
) (*model.ProductVariant, error) {
	validatedFilterFields, err := getProdVariantUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	productVariant, err := getProdVariantUcase.getProductVariantRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getProdVariantUcase.pathIdentity,
			err,
		)
	}
	return productVariant, nil
}
