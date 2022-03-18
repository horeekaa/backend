package mongodbpaymentdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbpaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewPaymentDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo, error) {
	basicOperation.SetCollection("payments")
	return &paymentDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "PaymentDataSource",
	}, nil
}

func (paymentDataSourceMongo *paymentDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (paymentDataSourceMongo *paymentDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error) {
	var output model.Payment
	_, err := paymentDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (paymentDataSourceMongo *paymentDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error) {
	var output model.Payment
	_, err := paymentDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (paymentDataSourceMongo *paymentDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Payment, error) {
	var payments = []*model.Payment{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var payment model.Payment
		if err := cursor.Decode(&payment); err != nil {
			return err
		}
		payments = append(payments, &payment)
		return nil
	}
	_, err := paymentDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return payments, err
}

func (paymentDataSourceMongo *paymentDataSourceMongo) Create(input *model.DatabaseCreatePayment, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error) {
	var outputModel model.Payment
	_, err := paymentDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (paymentDataSourceMongo *paymentDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdatePayment,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Payment, error) {
	existingObject, err := paymentDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			paymentDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.Payment
	_, err = paymentDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
