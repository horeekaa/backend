package googlecloudstoragecoreclientinterfaces

import (
	"cloud.google.com/go/storage"
)

type GoogleCloudStorageClient interface {
	Initialize() (bool, error)
	GetClient() (*storage.Client, error)
}
