package productdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetProductRepository interface {
	Execute(filterFields *model.ProductFilterFields) (*model.Product, error)
}
