package paymentdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePaymentTransactionComponent interface {
	PreTransaction(
		createPaymentInput *model.InternalCreatePayment,
	) (*model.InternalCreatePayment, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createPaymentInput *model.InternalCreatePayment,
	) (*model.Payment, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreatePaymentRepository interface {
	RunTransaction(
		createPaymentInput *model.InternalCreatePayment,
	) (*model.Payment, error)
}
