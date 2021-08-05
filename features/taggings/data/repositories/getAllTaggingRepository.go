package taggingdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingdomainrepositorytypes "github.com/horeekaa/backend/features/taggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllTaggingRepository struct {
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource
}

func NewGetAllTaggingRepository(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
) (taggingdomainrepositoryinterfaces.GetAllTaggingRepository, error) {
	return &getAllTaggingRepository{
		taggingDataSource,
	}, nil
}

func (getAllTaggingRepo *getAllTaggingRepository) Execute(
	input taggingdomainrepositorytypes.GetAllTaggingInput,
) ([]*model.Tagging, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	taggings, err := getAllTaggingRepo.taggingDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllTagging",
			err,
		)
	}

	return taggings, nil
}
