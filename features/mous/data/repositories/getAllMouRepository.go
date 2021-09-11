package moudomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositorytypes "github.com/horeekaa/backend/features/mous/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllMouRepository struct {
	mouDataSource databasemoudatasourceinterfaces.MouDataSource
}

func NewGetAllMouRepository(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
) (moudomainrepositoryinterfaces.GetAllMouRepository, error) {
	return &getAllMouRepository{
		mouDataSource,
	}, nil
}

func (getAllMmbAccRefRepo *getAllMouRepository) Execute(
	input moudomainrepositorytypes.GetAllMouInput,
) ([]*model.Mou, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	mous, err := getAllMmbAccRefRepo.mouDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllMou",
			err,
		)
	}

	return mous, nil
}
