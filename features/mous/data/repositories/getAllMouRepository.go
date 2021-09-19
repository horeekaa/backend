package moudomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositorytypes "github.com/horeekaa/backend/features/mous/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllMouRepository struct {
	mouDataSource     databasemoudatasourceinterfaces.MouDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllMouRepository(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (moudomainrepositoryinterfaces.GetAllMouRepository, error) {
	return &getAllMouRepository{
		mouDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllMmbAccRefRepo *getAllMouRepository) Execute(
	input moudomainrepositorytypes.GetAllMouInput,
) ([]*model.Mou, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMmbAccRefRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

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
