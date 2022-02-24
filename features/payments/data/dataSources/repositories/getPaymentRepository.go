package paymentdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getPaymentRepository struct {
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource
}

func NewGetPaymentRepository(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
) (paymentdomainrepositoryinterfaces.GetPaymentRepository, error) {
	return &getPaymentRepository{
		paymentDataSource,
	}, nil
}

func (getPaymentRepo *getPaymentRepository) Execute(filterFields *model.PaymentFilterFields) (*model.Payment, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	payment, err := getPaymentRepo.paymentDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getPayment",
			err,
		)
	}

	return payment, nil
}
