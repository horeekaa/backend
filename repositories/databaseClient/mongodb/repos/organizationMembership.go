package mongorepos

import (
	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperationinterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/operations"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type organizationMembershipRepoMongo struct {
	basicOperation mongooperationinterfaces.BasicOperation
}

func NeworganizationMembershipRepoMongo(mongoRepo *databaseclient.MongoRepository) (mongorepointerfaces.OrganizationMembershipRepoMongo, error) {
	basicOperation, err := mongooperations.NewBasicOperation(
		(*mongoRepo).Client,
		(*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("organizationmemberships"),
		(*mongoRepo).Timeout,
		"organizationmemberships",
	)
	if err != nil {
		return nil, err
	}

	return &organizationMembershipRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.OrganizationMembership
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.OrganizationMembership
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.OrganizationMembership, error) {
	var organizationMemberships = []*model.OrganizationMembership{}
	cursorDecoder := func(cursor *mongooperationmodels.CursorObject) (interface{}, error) {
		var organizationMembership *model.OrganizationMembership
		err := cursor.MongoFindCursor.Decode(organizationMembership)
		if err != nil {
			return nil, err
		}
		organizationMemberships = append(organizationMemberships, organizationMembership)
		return nil, nil
	}

	_, err := orgMemberRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return organizationMemberships, err
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) Create(input *model.CreateOrganizationMembership, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := orgMemberRepoMongo.basicOperation.Create(*defaultedInput.CreateOrganizationMembership, operationOptions)
	if err != nil {
		return nil, err
	}

	organizationMembershipOutput := output.Object.(model.OrganizationMembership)

	organizationMembership := &model.OrganizationMembership{
		ID:            output.ID,
		Person:        organizationMembershipOutput.Person,
		Role:          organizationMembershipOutput.Role,
		Status:        organizationMembershipOutput.Status,
		Access:        organizationMembershipOutput.Access,
		DefaultAccess: organizationMembershipOutput.DefaultAccess,
	}

	return organizationMembership, err
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) Update(ID interface{}, updateData *model.UpdateOrganizationMembership, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := orgMemberRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdateOrganizationMembership, operationOptions)
	var output model.OrganizationMembership
	res.Decode(&output)

	return &output, err
}

type setOrganizationMembershipDefaultValuesOutput struct {
	CreateOrganizationMembership *model.CreateOrganizationMembership
	UpdateOrganizationMembership *model.UpdateOrganizationMembership
}

func (orgMemberRepoMongo *organizationMembershipRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperationmodels.OperationOptions) (*setOrganizationMembershipDefaultValuesOutput, error) {
	updateInput := input.(model.UpdateOrganizationMembership)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
		_, err := orgMemberRepoMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		return &setOrganizationMembershipDefaultValuesOutput{
			UpdateOrganizationMembership: &model.UpdateOrganizationMembership{
				ID:     updateInput.ID,
				Person: updateInput.Person,
				Role:   updateInput.Role,
				Status: updateInput.Status,
				Access: updateInput.Access,
			},
		}, nil
	}
	createInput := (input).(model.CreateOrganizationMembership)

	return &setOrganizationMembershipDefaultValuesOutput{
		CreateOrganizationMembership: &model.CreateOrganizationMembership{
			Person: createInput.Person,
			Role:   createInput.Role,
			Status: createInput.Status,
			Access: createInput.Access,
		},
	}, nil
}
