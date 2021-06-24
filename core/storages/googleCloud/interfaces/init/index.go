package googlecloudstoragecoreclientinterfaces

import (
	"io"

	googlecloudstoragecorewrapperinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/wrappers"
)

type GoogleCloudStorageClient interface {
	Initialize() (bool, error)
	CopyWrite(dst io.Writer, src io.Reader) (written int64, err error)
	GetObjectHandle(bucketName string, objectPath string) (googlecloudstoragecorewrapperinterfaces.GCSObjectHandle, error)
}
