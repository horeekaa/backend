package mongodbaddressdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type addressDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewAddressDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaddressdatasourceinterfaces.AddressDataSourceMongo, error) {
	basicOperation.SetCollection("addresses")
	return &addressDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "AddressDataSource",
	}, nil
}

func (addrDataSourceMongo *addressDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (addrDataSourceMongo *addressDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error) {
	var output model.Address
	_, err := addrDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (addrDataSourceMongo *addressDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error) {
	var output model.Address
	_, err := addrDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (addrDataSourceMongo *addressDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Address, error) {
	var addresss = []*model.Address{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var address model.Address
		if err := cursor.Decode(&address); err != nil {
			return err
		}
		addresss = append(addresss, &address)
		return nil
	}
	_, err := addrDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return addresss, err
}

func (addrDataSourceMongo *addressDataSourceMongo) Create(input *model.DatabaseCreateAddress, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error) {
	var outputModel model.Address
	_, err := addrDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (addrDataSourceMongo *addressDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateAddress, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error) {
	existingObject, err := addrDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			addrDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.Address
	_, err = addrDataSourceMongo.basicOperation.Update(
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
