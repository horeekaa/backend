package mongodbmemberaccessrefdatasources

import (
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
	pathIdentity   string
}

func NewMemberAccessRefDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo, error) {
	basicOperation.SetCollection("memberaccessrefs")
	return &memberAccessRefDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "MemberAccessRefDataSource",
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

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) Create(input *model.DatabaseCreateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error) {
	var outputModel model.MemberAccessRef
	_, err := mmbAccRefDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mmbAccRefDataSourceMongo *memberAccessRefDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateMemberAccessRef,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.MemberAccessRef, error) {
	existingObject, err := mmbAccRefDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			mmbAccRefDataSourceMongo.pathIdentity,
			nil,
		)
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
