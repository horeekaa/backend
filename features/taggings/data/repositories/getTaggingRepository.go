package taggingdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getTaggingRepository struct {
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource
}

func NewGetTaggingRepository(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
) (taggingdomainrepositoryinterfaces.GetTaggingRepository, error) {
	return &getTaggingRepository{
		taggingDataSource,
	}, nil
}

func (getTaggingRepo *getTaggingRepository) Execute(filterFields *model.TaggingFilterFields) (*model.Tagging, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	tagging, err := getTaggingRepo.taggingDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getTagging",
			err,
		)
	}

	return tagging, nil
}
