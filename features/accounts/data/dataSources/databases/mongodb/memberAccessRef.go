package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
)

type memberAccessRefRepoMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMemberAccessRefRepoMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo, error) {
	basicOperation.SetCollection("memberaccessrefs")
	return &memberAccessRefRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.MemberAccessRef
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.MemberAccessRef
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.MemberAccessRef, error) {
	var memberAccessRefs = []*model.MemberAccessRef{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var memberAccessRef *model.MemberAccessRef
		err := cursor.MongoFindCursor.Decode(memberAccessRef)
		if err != nil {
			return nil, err
		}
		memberAccessRefs = append(memberAccessRefs, memberAccessRef)
		return nil, nil
	}

	_, err := orgMemberRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccessRefs, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Create(input *model.CreateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := orgMemberRepoMongo.basicOperation.Create(*defaultedInput.CreateMemberAccessRef, operationOptions)
	if err != nil {
		return nil, err
	}

	memberAccessRefOutput := output.Object.(model.MemberAccessRef)

	memberAccessRef := &model.MemberAccessRef{
		ID:                         output.ID,
		Access:                     memberAccessRefOutput.Access,
		MemberAccessRefType:        memberAccessRefOutput.MemberAccessRefType,
		OrganizationMembershipRole: memberAccessRefOutput.OrganizationMembershipRole,
		OrganizationType:           memberAccessRefOutput.OrganizationType,
		CreatedAt:                  memberAccessRefOutput.CreatedAt,
		UpdatedAt:                  memberAccessRefOutput.UpdatedAt,
	}

	return memberAccessRef, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Update(ID interface{}, updateData *model.UpdateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := orgMemberRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdateMemberAccessRef, operationOptions)
	var output model.MemberAccessRef
	res.Decode(&output)

	return &output, err
}

type setMemberAccessRefDefaultValuesOutput struct {
	CreateMemberAccessRef *model.CreateMemberAccessRef
	UpdateMemberAccessRef *model.UpdateMemberAccessRef
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMemberAccessRefDefaultValuesOutput, error) {
	var currentTime = time.Now()

	updateInput := input.(model.UpdateMemberAccessRef)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		_, err := orgMemberRepoMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		return &setMemberAccessRefDefaultValuesOutput{
			UpdateMemberAccessRef: &model.UpdateMemberAccessRef{
				ID:                         updateInput.ID,
				Access:                     updateInput.Access,
				MemberAccessRefType:        updateInput.MemberAccessRefType,
				OrganizationMembershipRole: updateInput.OrganizationMembershipRole,
				OrganizationType:           updateInput.OrganizationType,
				UpdatedAt:                  &currentTime,
			},
		}, nil
	}
	createInput := (input).(model.CreateMemberAccessRef)

	return &setMemberAccessRefDefaultValuesOutput{
		CreateMemberAccessRef: &model.CreateMemberAccessRef{
			Access:                     createInput.Access,
			MemberAccessRefType:        createInput.MemberAccessRefType,
			OrganizationMembershipRole: createInput.OrganizationMembershipRole,
			OrganizationType:           createInput.OrganizationType,
			CreatedAt:                  &currentTime,
			UpdatedAt:                  &currentTime,
		},
	}, nil
}
