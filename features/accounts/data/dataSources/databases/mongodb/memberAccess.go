package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
)

type memberAccessRepoMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMemberAccessRepoMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo, error) {
	basicOperation.SetCollection("memberaccesses")
	return &memberAccessRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (memberAccRepoMongo *memberAccessRepoMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	res, err := memberAccRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.MemberAccess
	res.Decode(&output)
	return &output, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.MemberAccess, error) {
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

	_, err := memberAccRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccesses, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Create(input *model.CreateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccRepoMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
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
		ID:                         output.ID,
		Account:                    memberAccessOutput.Account,
		Organization:               memberAccessOutput.Organization,
		OrganizationMembershipRole: memberAccessOutput.OrganizationMembershipRole,
		MemberAccessRefType:        memberAccessOutput.MemberAccessRefType,
		Access:                     memberAccessOutput.Access,
		DefaultAccess:              memberAccessOutput.DefaultAccess,
		Status:                     memberAccessOutput.Status,
		CreatedAt:                  memberAccessOutput.CreatedAt,
		UpdatedAt:                  memberAccessOutput.UpdatedAt,
	}

	return memberAccess, err
}

func (memberAccRepoMongo *memberAccessRepoMongo) Update(ID interface{}, updateData *model.UpdateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	defaultedInput, err := memberAccRepoMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
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

func (memberAccRepoMongo *memberAccessRepoMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMemberAccessDefaultValuesOutput, error) {
	var currentTime = time.Now()

	updateInput := input.(model.UpdateMemberAccess)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		_, err := memberAccRepoMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		return &setMemberAccessDefaultValuesOutput{
			UpdateMemberAccess: &model.UpdateMemberAccess{
				ID:                         updateInput.ID,
				OrganizationMembershipRole: updateInput.OrganizationMembershipRole,
				Access:                     updateInput.Access,
				Status:                     updateInput.Status,
				UpdatedAt:                  &currentTime,
			},
		}, nil
	}
	createInput := (input).(model.CreateMemberAccess)

	return &setMemberAccessDefaultValuesOutput{
		CreateMemberAccess: &model.CreateMemberAccess{
			Account:                    createInput.Account,
			Organization:               createInput.Organization,
			OrganizationMembershipRole: createInput.OrganizationMembershipRole,
			MemberAccessRefType:        createInput.MemberAccessRefType,
			Access:                     createInput.Access,
			Status:                     createInput.Status,
			CreatedAt:                  &currentTime,
			UpdatedAt:                  &currentTime,
		},
	}, nil
}
