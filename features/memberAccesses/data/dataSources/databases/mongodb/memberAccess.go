package mongodbmemberaccessdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbmemberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type memberAccessDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewMemberAccessDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo, error) {
	basicOperation.SetCollection("memberaccesses")
	return &memberAccessDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "MemberAccessDataSource",
	}, nil
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	var output model.MemberAccess
	_, err := memberAccDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	var output model.MemberAccess
	_, err := memberAccDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.MemberAccess, error) {
	var memberAccesses = []*model.MemberAccess{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var memberAccess model.MemberAccess
		if err := cursor.Decode(&memberAccess); err != nil {
			return err
		}
		memberAccesses = append(memberAccesses, &memberAccess)
		return nil
	}
	_, err := memberAccDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return memberAccesses, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Create(input *model.DatabaseCreateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error) {
	var outputModel model.MemberAccess
	_, err := memberAccDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (memberAccDataSourceMongo *memberAccessDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateMemberAccess,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.MemberAccess, error) {
	existingObject, err := memberAccDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			memberAccDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.MemberAccess
	_, err = memberAccDataSourceMongo.basicOperation.Update(
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
