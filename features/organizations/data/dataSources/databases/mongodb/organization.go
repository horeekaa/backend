package mongodborganizationdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type organizationDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewOrganizationDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo, error) {
	basicOperation.SetCollection("organizations")
	return &organizationDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "OrganizationDataSource",
	}, nil
}

func (orgDataSourceMongo *organizationDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (orgDataSourceMongo *organizationDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	var output model.Organization
	_, err := orgDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (orgDataSourceMongo *organizationDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	var output model.Organization
	_, err := orgDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Organization, error) {
	var organizations = []*model.Organization{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var organization model.Organization
		if err := cursor.Decode(&organization); err != nil {
			return err
		}
		organizations = append(organizations, &organization)
		return nil
	}
	_, err := orgDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return organizations, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Create(input *model.DatabaseCreateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	_, err := orgDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Organization
	_, err = orgDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateOrganization,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Organization, error) {
	_, err := orgDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Organization
	_, err = orgDataSourceMongo.basicOperation.Update(
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

func (orgDataSourceMongo *organizationDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateOrganization,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := orgDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			orgDataSourceMongo.pathIdentity,
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (orgDataSourceMongo *organizationDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateOrganization,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultPoint := 0

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}
	if input.ProfilePhotos == nil {
		input.ProfilePhotos = []*model.ObjectIDOnly{}
	}
	if input.Taggings == nil {
		input.Taggings = []*model.ObjectIDOnly{}
	}
	input.Point = &defaultPoint
	input.UnfinalizedPoint = &defaultPoint
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
