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

type memberAccessRefRepoMongo struct {
	basicOperation mongooperationinterfaces.BasicOperation
}

func NewMemberAccessRefRepoMongo(mongoRepo *databaseclient.MongoRepository) (mongorepointerfaces.MemberAccessRefRepoMongo, error) {
	basicOperation, err := mongooperations.NewBasicOperation(
		(*mongoRepo).Client,
		(*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("memberaccessrefs"),
		(*mongoRepo).Timeout,
		"memberaccessrefs",
	)
	if err != nil {
		return nil, err
	}

	return &memberAccessRefRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.MemberAccessRef
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error) {
	res, err := orgMemberRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.MemberAccessRef
	res.Decode(&output)
	return &output, err
}

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.MemberAccessRef, error) {
	var memberAccessRefs = []*model.MemberAccessRef{}
	cursorDecoder := func(cursor *mongooperationmodels.CursorObject) (interface{}, error) {
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

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Create(input *model.CreateMemberAccessRef, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
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

func (orgMemberRepoMongo *memberAccessRefRepoMongo) Update(ID interface{}, updateData *model.UpdateMemberAccessRef, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error) {
	defaultedInput, err := orgMemberRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
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

func (orgMemberRepoMongo *memberAccessRefRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperationmodels.OperationOptions) (*setMemberAccessRefDefaultValuesOutput, error) {
	var currentTime = time.Now()

	updateInput := input.(model.UpdateMemberAccessRef)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
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
