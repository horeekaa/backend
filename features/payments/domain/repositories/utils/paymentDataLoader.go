package paymentdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type PaymentLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		invoice *model.InvoiceForPaymentInput,
		organization *model.OrganizationForPaymentInput,
	) (bool, error)
}
