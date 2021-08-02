package tagdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagdomainrepositorytypes "github.com/horeekaa/backend/features/tags/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllTagRepository struct {
	tagDataSource databasetagdatasourceinterfaces.TagDataSource
}

func NewGetAllTagRepository(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
) (tagdomainrepositoryinterfaces.GetAllTagRepository, error) {
	return &getAllTagRepository{
		tagDataSource,
	}, nil
}

func (getAllTagRepo *getAllTagRepository) Execute(
	input tagdomainrepositorytypes.GetAllTagInput,
) ([]*model.Tag, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

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
