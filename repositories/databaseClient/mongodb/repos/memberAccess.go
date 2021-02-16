package mongorepos

import (
	"time"

	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperationinterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/operations"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type memberAccessRepoMongo struct {
	basicOperation mongooperationinterfaces.BasicOperation
}

func NewMemberAccessRepoMongo(mongoRepo *databaseclient.MongoRepository) (mongorepointerfaces.MemberAccessRepoMongo, error) {
	basicOperation, err := mongooperations.NewBasicOperation(
		(*mongoRepo).Client,
		(*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("MemberAccesses"),
		(*mongoRepo).Timeout,
		"memberaccesses",
	)
	if err != nil {
		return nil, err
	}

	return &memberAccessRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (memberAccRepoMongo *memberAccessRepoMongo) FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.MemberAccess, error) {
	var memberAccesses = []*model.MemberAccess{}
	cursorDecoder := func(cursor *mongooperationmodels.CursorObject) (interface{}, error) {
		var memberAccess *model.MemberAccess
		err := cursor.MongoFindCursor.Decode(memberAccess)
		if err != nil {
			return nil, err
		}
		memberAccesses = append(memberAccesses, memberAccess)
		return nil, nil
	}

	_, err := memberAccRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccesses, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Create(input *model.CreateMemberAccess, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := memberAccRepoMongo.basicOperation.Create(*defaultedInput.CreateMemberAccess, operationOptions)
	if err != nil {
		return nil, err
	}

	memberAccessOutput := output.Object.(model.MemberAccess)

	memberAccess := &model.MemberAccess{
		ID:               output.ID,
		Access:           memberAccessOutput.Access,
		Role:             memberAccessOutput.Role,
		OrganizationType: memberAccessOutput.OrganizationType,
		CreatedAt:        memberAccessOutput.CreatedAt,
		UpdatedAt:        memberAccessOutput.UpdatedAt,
	}

	return memberAccess, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Update(ID interface{}, updateData *model.UpdateMemberAccess, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := memberAccRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdateMemberAccess, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)

	return &output, err
}

type setMemberAccessDefaultValuesOutput struct {
	CreateMemberAccess *model.CreateMemberAccess
	UpdateMemberAccess *model.UpdateMemberAccess
}

func (memberAccRepoMongo *memberAccessRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperationmodels.OperationOptions) (*setMemberAccessDefaultValuesOutput, error) {
	var currentTime = time.Now()

	updateInput := input.(model.UpdateMemberAccess)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
		_, err := memberAccRepoMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		return &setMemberAccessDefaultValuesOutput{
			UpdateMemberAccess: &model.UpdateMemberAccess{
				ID:               updateInput.ID,
				Access:           updateInput.Access,
				Role:             updateInput.Role,
				OrganizationType: updateInput.OrganizationType,
				UpdatedAt:        &currentTime,
			},
		}, nil
	}
	createInput := (input).(model.CreateMemberAccess)

	return &setMemberAccessDefaultValuesOutput{
		CreateMemberAccess: &model.CreateMemberAccess{
			Access:           createInput.Access,
			Role:             createInput.Role,
			OrganizationType: createInput.OrganizationType,
			CreatedAt:        &currentTime,
			UpdatedAt:        &currentTime,
		},
	}, nil
}
