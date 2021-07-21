package descriptivephotodomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateDescriptivePhotoTransactionComponent interface {
	PreTransaction(
		updateDescriptivePhotoInput *model.InternalUpdateDescriptivePhoto,
	) (*model.InternalUpdateDescriptivePhoto, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateDescriptivePhotoInput *model.InternalUpdateDescriptivePhoto,
	) (*model.DescriptivePhoto, error)
}
