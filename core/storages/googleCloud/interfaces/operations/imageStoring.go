package googlecloudstoragecoreoperationinterfaces

import (
	"context"

	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
	"github.com/horeekaa/backend/model"
)

type GCSBasicImageStoringOperation interface {
	UploadImage(
		ctx context.Context,
		category model.DescriptivePhotoCategory,
		file googlecloudstoragecoretypes.GCSFileUpload,
	) (string, error)
	DeleteImage(
		ctx context.Context,
		url string,
	) (bool, error)
}
