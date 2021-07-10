package productdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllProductInput struct {
	FilterFields  *model.ProductFilterFields
	PaginationOpt *model.PaginationOptionInput
}
