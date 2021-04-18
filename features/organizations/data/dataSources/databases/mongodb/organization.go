package mongodborganizationdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
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

func (orgDataSourceMongo *organizationDataSourceMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	res, err := orgDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Organization
	res.Decode(&output)
	return &output, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	res, err := orgDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Organization
	res.Decode(&output)
	return &output, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Organization, error) {
	var organizations = []*model.Organization{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var organization *model.Organization
		err := cursor.MongoFindCursor.Decode(organization)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
		return nil, nil
	}

	_, err := orgDataSourceMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
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

	organizationOutput := output.Object.(model.Organization)
	organizationOutput.ID = output.ID

	return &organizationOutput, err
}

func (orgDataSourceMongo *organizationDataSourceMongo) Update(ID interface{}, updateData *model.UpdateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error) {
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := orgDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateOrganization, operationOptions)
	var output model.Organization
	res.Decode(&output)

	return &output, err
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