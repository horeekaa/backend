package tagdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getTagRepository struct {
	tagDataSource databasetagdatasourceinterfaces.TagDataSource
	pathIdentity  string
}

func NewGetTagRepository(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
) (tagdomainrepositoryinterfaces.GetTagRepository, error) {
	return &getTagRepository{
		tagDataSource,
		"GetTagRepository",
	}, nil
}

func (getTagRepo *getTagRepository) Execute(filterFields *model.TagFilterFields) (*model.Tag, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	tag, err := getTagRepo.tagDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getTagRepo.pathIdentity,
			err,
		)
	}

	return tag, nil
}
