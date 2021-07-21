package descriptivephotodomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateDescriptivePhotoTransactionComponent interface {
	PreTransaction(
		createDescriptivePhotoInput *model.InternalCreateDescriptivePhoto,
	) (*model.InternalCreateDescriptivePhoto, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createDescriptivePhotoInput *model.InternalCreateDescriptivePhoto,
	) (*model.DescriptivePhoto, error)
}
