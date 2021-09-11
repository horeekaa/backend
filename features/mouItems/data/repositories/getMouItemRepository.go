package mouitemdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getMouItemRepository struct {
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource
}

func NewGetMouItemRepository(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
) (mouitemdomainrepositoryinterfaces.GetMouItemRepository, error) {
	return &getMouItemRepository{
		mouItemDataSource,
	}, nil
}

func (getMouItemRefRepo *getMouItemRepository) Execute(filterFields *model.MouItemFilterFields) (*model.MouItem, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mouItem, err := getMouItemRefRepo.mouItemDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getMouItem",
			err,
		)
	}

	return mouItem, nil
}
