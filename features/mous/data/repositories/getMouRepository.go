package moudomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getMouRepository struct {
	mouDataSource databasemoudatasourceinterfaces.MouDataSource
}

func NewGetMouRepository(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
) (moudomainrepositoryinterfaces.GetMouRepository, error) {
	return &getMouRepository{
		mouDataSource,
	}, nil
}

func (getOrgRepo *getMouRepository) Execute(filterFields *model.MouFilterFields) (*model.Mou, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mou, err := getOrgRepo.mouDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getMou",
			err,
		)
	}

	return mou, nil
}
