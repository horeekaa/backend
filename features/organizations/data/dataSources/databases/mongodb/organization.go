package mongodborganizationdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type organizationDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewOrganizationDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo, error) {
	basicOperation.SetCollection("organizations")
	return &organizationDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
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
	_, err := orgDataSourceMongo.basicOperation.Find(query, paginationOpts, &organizations, operationOptions)
	if err != nil {
		return nil, err
	}

	return organizations, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Create(input *model.CreateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Organization
	_, err = orgDataSourceMongo.basicOperation.Create(*defaultedInput.CreateOrganization, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	updateData.ID = ID
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Organization
	_, err = orgDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateOrganization, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setorganizationDefaultValuesOutput struct {
	CreateOrganization *model.CreateOrganization
	UpdateOrganization *model.UpdateOrganization
}

func (orgDataSourceMongo *organizationDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setorganizationDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultPoint := 0

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.UpdateOrganization)
		_, err := orgDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setorganizationDefaultValuesOutput{
			UpdateOrganization: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreateOrganization)
	if createInput.ProposalStatus == nil {
		createInput.ProposalStatus = &defaultProposalStatus
	}
	createInput.Point = &defaultPoint
	createInput.UnfinalizedPoint = &defaultPoint
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setorganizationDefaultValuesOutput{
		CreateOrganization: &createInput,
	}, nil
}
