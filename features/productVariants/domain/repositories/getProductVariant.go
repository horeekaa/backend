package productvariantdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetProductVariantRepository interface {
	Execute(filterFields *model.ProductVariantFilterFields) (*model.ProductVariant, error)
}
