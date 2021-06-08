package googlecloudstoragecoreoperationinterfaces

import (
	"context"

	storageenums "github.com/horeekaa/backend/core/storages/enums"
	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
)

type GSCBasicImageStoringOperation interface {
	UploadImage(
		ctx context.Context,
		category storageenums.StorageCategory,
		file googlecloudstoragecoretypes.GCSFileUpload,
	) (string, error)
	DeleteImage(
		ctx context.Context,
		url string,
	) (bool, error)
}
