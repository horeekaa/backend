package purchaseorderitempresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type UpdatePurchaseOrderItemDeliveryUsecaseInput struct {
	Context                         context.Context
	UpdatePurchaseOrderItemDelivery *model.UpdatePurchaseOrderItemDeliveryOnly
}
