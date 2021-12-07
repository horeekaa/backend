package taggingdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	taggingdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories/utils"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type taggingLoader struct {
	tagDataSource databasetagdatasourceinterfaces.TagDataSource
}

func NewTaggingLoader(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
) (taggingdomainrepositoryutilityinterfaces.TaggingLoader, error) {
	return &taggingLoader{
		tagDataSource: tagDataSource,
	}, nil
}

func (taggingLoaderUtil *taggingLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	tagForTagging *model.TagForTaggingInput,
) (bool, error) {
	if tagForTagging == nil {
		return true, nil
	}
	loadedTag, err := taggingLoaderUtil.tagDataSource.GetMongoDataSource().FindByID(
		*tagForTagging.ID,
		session,
	)
	if err != nil {
		return false, err
	}

	jsonTemp, _ := json.Marshal(loadedTag)
	json.Unmarshal(jsonTemp, &tagForTagging)

	return true, nil
}
