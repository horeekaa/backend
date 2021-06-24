package googlecloudstoragecoreclientinterfaces

import (
	googlecloudstoragecorewrapperinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/wrappers"
)

type GoogleCloudStorageClient interface {
	Initialize() (bool, error)
	GetObjectHandle(bucketName string, objectPath string) (googlecloudstoragecorewrapperinterfaces.GCSObjectHandle, error)
}
