package mongodbmemberaccessrefdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
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

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	var output model.MemberAccessRef
	_, err := mmbAccRefDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	var output model.MemberAccessRef
	_, err := mmbAccRefDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpt *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.MemberAccessRef, error) {
	var memberAccessRefs = []*model.MemberAccessRef{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var memberAccessRef model.MemberAccessRef
		if err := cursor.Decode(&memberAccessRef); err != nil {
			return err
		}
		memberAccessRefs = append(memberAccessRefs, &memberAccessRef)
		return nil
	}
	_, err := mmbAccRefDataSourceMongo.basicOperation.Find(query, paginationOpt, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccessRefs, err
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) Create(input *model.InternalCreateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	_, err := mmbAccRefDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.MemberAccessRef
	_, err = mmbAccRefDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.InternalUpdateMemberAccessRef,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.MemberAccessRef, error) {
	_, err := mmbAccRefDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.MemberAccessRef
	_, err = mmbAccRefDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.InternalUpdateMemberAccessRef,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	var currentTime = time.Now()
	existingObject, err := mmbAccRefDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/memberAccessRefDataSource/update",
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.InternalCreateMemberAccessRef,
) (bool, error) {
	var currentTime = time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
