package taggingdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingdomainrepositorytypes "github.com/horeekaa/backend/features/taggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllTaggingRepository struct {
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity      string
}

func NewGetAllTaggingRepository(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (taggingdomainrepositoryinterfaces.GetAllTaggingRepository, error) {
	return &getAllTaggingRepository{
		taggingDataSource,
		mongoQueryBuilder,
		"GetAllTaggingRepository",
	}, nil
}

func (getAllTaggingRepo *getAllTaggingRepository) Execute(
	input taggingdomainrepositorytypes.GetAllTaggingInput,
) ([]*model.Tagging, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllTaggingRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	taggings, err := getAllTaggingRepo.taggingDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllTaggingRepo.pathIdentity,
			err,
		)
	}

	return taggings, nil
}
