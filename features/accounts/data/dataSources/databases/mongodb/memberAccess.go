package mongodbaccountdatasources

import (
	"encoding/json"
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccess
	res.Decode(&output)
	return &output, nil
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccess
	err = res.Decode(&output)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &output, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.MemberAccess, error) {
	var memberAccesses = []*model.MemberAccess{}
	cursorDecoder := func(cursor *mongo.Cursor) (interface{}, error) {
		var memberAccess model.MemberAccess
		err := cursor.Decode(&memberAccess)
		if err != nil {
			return nil, err
		}
		memberAccesses = append(memberAccesses, &memberAccess)
		return nil, nil
	}

	_, err := memberAccDataSourceMongo.basicOperation.Find(query, paginationOpts, cursorDecoder, operationOptions)
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

	var outputModel model.MemberAccess
	jsonTemp, _ := json.Marshal(output.Object)
	json.Unmarshal(jsonTemp, &outputModel)
	outputModel.ID = output.ID

	return &outputModel, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	updateData.ID = ID
	defaultedInput, err := memberAccDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := memberAccDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateMemberAccess, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccess
	res.Decode(&output)

	return &output, nil
}

type setMemberAccessDefaultValuesOutput struct {
	CreateMemberAccess *model.CreateMemberAccess
	UpdateMemberAccess *model.UpdateMemberAccess
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMemberAccessDefaultValuesOutput, error) {
	var currentTime = time.Now()

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.UpdateMemberAccess)
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
