package mouitempresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type GetAllMouItemUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.MouItemFilterFields
	PaginationOps *model.PaginationOptionInput
}
