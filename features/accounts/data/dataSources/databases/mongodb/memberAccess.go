package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
)

type memberAccessDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMemberAccessDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo, error) {
	basicOperation.SetCollection("memberaccesses")
	return &memberAccessDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.MemberAccess, error) {
	var memberAccesses = []*model.MemberAccess{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var memberAccess *model.MemberAccess
		err := cursor.MongoFindCursor.Decode(memberAccess)
		if err != nil {
			return nil, err
		}
		memberAccesses = append(memberAccesses, memberAccess)
		return nil, nil
	}

	_, err := memberAccDataSourceMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccesses, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Create(input *model.CreateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := memberAccDataSourceMongo.basicOperation.Create(*defaultedInput.CreateMemberAccess, operationOptions)
	if err != nil {
		return nil, err
	}

	memberAccessOutput := output.Object.(model.MemberAccess)
	memberAccessOutput.ID = output.ID

	return &memberAccessOutput, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Update(ID interface{}, updateData *model.UpdateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := memberAccDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateMemberAccess, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)

	return &output, err
}

type setMemberAccessDefaultValuesOutput struct {
	CreateMemberAccess *model.CreateMemberAccess
	UpdateMemberAccess *model.UpdateMemberAccess
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMemberAccessDefaultValuesOutput, error) {
	var currentTime = time.Now()

	updateInput := input.(model.UpdateMemberAccess)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		_, err := memberAccDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setMemberAccessDefaultValuesOutput{
			UpdateMemberAccess: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreateMemberAccess)
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setMemberAccessDefaultValuesOutput{
		CreateMemberAccess: &createInput,
	}, nil
}
