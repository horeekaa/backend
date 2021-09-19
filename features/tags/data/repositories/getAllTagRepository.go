package tagdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagdomainrepositorytypes "github.com/horeekaa/backend/features/tags/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllTagRepository struct {
	tagDataSource     databasetagdatasourceinterfaces.TagDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllTagRepository(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (tagdomainrepositoryinterfaces.GetAllTagRepository, error) {
	return &getAllTagRepository{
		tagDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllTagRepo *getAllTagRepository) Execute(
	input tagdomainrepositorytypes.GetAllTagInput,
) ([]*model.Tag, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllTagRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	tags, err := getAllTagRepo.tagDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllTag",
			err,
		)
	}

	return tags, nil
}
