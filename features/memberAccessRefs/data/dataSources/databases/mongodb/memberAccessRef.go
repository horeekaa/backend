package mongodbmemberaccessrefdatasources

import (
	"encoding/json"
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbmemberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberAccessRefDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMemberAccessRefDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo, error) {
	basicOperation.SetCollection("memberaccessrefs")
	return &memberAccessRefDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccessRef
	res.Decode(&output)
	return &output, nil
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccessRef
	err = res.Decode(&output)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &output, err
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpt *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.MemberAccessRef, error) {
	var memberAccessRefs = []*model.MemberAccessRef{}
	cursorDecoder := func(cursor *mongo.Cursor) (interface{}, error) {
		var memberAccessRef model.MemberAccessRef
		err := cursor.Decode(&memberAccessRef)
		if err != nil {
			return nil, err
		}
		memberAccessRefs = append(memberAccessRefs, &memberAccessRef)
		return nil, nil
	}

	_, err := orgMemberDataSourceMongo.basicOperation.Find(query, paginationOpt, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccessRefs, err
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) Create(input *model.CreateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	defaultedInput, err := orgMemberDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := orgMemberDataSourceMongo.basicOperation.Create(*defaultedInput.CreateMemberAccessRef, operationOptions)
	if err != nil {
		return nil, err
	}

	var outputModel model.MemberAccessRef
	jsonTemp, _ := json.Marshal(output.Object)
	json.Unmarshal(jsonTemp, &outputModel)
	outputModel.ID = output.ID

	return &outputModel, err
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	updateData.ID = ID
	defaultedInput, err := orgMemberDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := orgMemberDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateMemberAccessRef, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccessRef
	res.Decode(&output)

	return &output, nil
}

type setMemberAccessRefDefaultValuesOutput struct {
	CreateMemberAccessRef *model.CreateMemberAccessRef
	UpdateMemberAccessRef *model.UpdateMemberAccessRef
}

func (orgMemberDataSourceMongo *memberAccessRefDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMemberAccessRefDefaultValuesOutput, error) {
	var currentTime = time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.UpdateMemberAccessRef)
		_, err := orgMemberDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setMemberAccessRefDefaultValuesOutput{
			UpdateMemberAccessRef: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreateMemberAccessRef)
	if createInput.ProposalStatus == nil {
		createInput.ProposalStatus = &defaultProposalStatus
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setMemberAccessRefDefaultValuesOutput{
		CreateMemberAccessRef: &createInput,
	}, nil
}
