package mongodborganizationdatasources

import (
	"encoding/json"
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
	res, err := orgDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Organization
	res.Decode(&output)
	return &output, nil
}

func (orgDataSourceMongo *organizationDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	res, err := orgDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Organization
	err = res.Decode(&output)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &output, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Organization, error) {
	var organizations = []*model.Organization{}
	cursorDecoder := func(cursor *mongo.Cursor) (interface{}, error) {
		var organization model.Organization
		err := cursor.Decode(&organization)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, &organization)
		return nil, nil
	}

	_, err := orgDataSourceMongo.basicOperation.Find(query, paginationOpts, cursorDecoder, operationOptions)
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

	output, err := orgDataSourceMongo.basicOperation.Create(*defaultedInput.CreateOrganization, operationOptions)
	if err != nil {
		return nil, err
	}

	var outputModel model.Organization
	jsonTemp, _ := json.Marshal(output.Object)
	json.Unmarshal(jsonTemp, &outputModel)
	outputModel.ID = output.ID

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

	res, err := orgDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateOrganization, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Organization
	res.Decode(&output)

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
