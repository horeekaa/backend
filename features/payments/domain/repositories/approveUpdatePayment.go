package paymentdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdatePaymentTransactionComponent interface {
	PreTransaction(
		updatePaymentInput *model.InternalUpdatePayment,
	) (*model.InternalUpdatePayment, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updatePaymentInput *model.InternalUpdatePayment,
	) (*model.Payment, error)
}

type ApproveUpdatePaymentRepository interface {
	RunTransaction(
		updatePaymentInput *model.InternalUpdatePayment,
	) (*model.Payment, error)
}
