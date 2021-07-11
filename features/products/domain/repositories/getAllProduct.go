package productdomainrepositoryinterfaces

import (
	productdomainrepositorytypes "github.com/horeekaa/backend/features/products/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllProductRepository interface {
	Execute(filterFields productdomainrepositorytypes.GetAllProductInput) ([]*model.Product, error)
}
