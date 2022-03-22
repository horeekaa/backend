package mouitemdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositorytypes "github.com/horeekaa/backend/features/mouItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllMouItemRepository struct {
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity      string
}

func NewGetAllMouItemRepository(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (mouitemdomainrepositoryinterfaces.GetAllMouItemRepository, error) {
	return &getAllMouItemRepository{
		mouItemDataSource,
		mongoQueryBuilder,
		"GetAllMouItemRepository",
	}, nil
}

func (getAllMouItemRepo *getAllMouItemRepository) Execute(
	input mouitemdomainrepositorytypes.GetAllMouItemInput,
) ([]*model.MouItem, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMouItemRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	mouItems, err := getAllMouItemRepo.mouItemDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllMouItemRepo.pathIdentity,
			err,
		)
	}

	return mouItems, nil
}
