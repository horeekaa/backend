package mongodbaddressRegionGroupdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type addressRegionGroupDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewAddressRegionGroupDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSourceMongo, error) {
	basicOperation.SetCollection("addressregiongroups")
	return &addressRegionGroupDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error) {
	var output model.AddressRegionGroup
	_, err := addrGroupDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error) {
	var output model.AddressRegionGroup
	_, err := addrGroupDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.AddressRegionGroup, error) {
	var addressRegionGroups = []*model.AddressRegionGroup{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var addressRegionGroup model.AddressRegionGroup
		if err := cursor.Decode(&addressRegionGroup); err != nil {
			return err
		}
		addressRegionGroups = append(addressRegionGroups, &addressRegionGroup)
		return nil
	}
	_, err := addrGroupDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return addressRegionGroups, err
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) Create(input *model.DatabaseCreateAddressRegionGroup, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error) {
	_, err := addrGroupDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.AddressRegionGroup
	_, err = addrGroupDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateAddressRegionGroup, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error) {
	_, err := addrGroupDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.AddressRegionGroup
	_, err = addrGroupDataSourceMongo.basicOperation.Update(
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

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateAddressRegionGroup,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()

	existingObject, err := addrGroupDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/addressRegionGroupDataSource/update",
			nil,
		)
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (addrGroupDataSourceMongo *addressRegionGroupDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateAddressRegionGroup,
) (bool, error) {
	currentTime := time.Now()

	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
