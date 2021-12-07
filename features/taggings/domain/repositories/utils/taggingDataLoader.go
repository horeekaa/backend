package taggingdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type TaggingLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		tagForTagging *model.TagForTaggingInput,
	) (bool, error)
}
