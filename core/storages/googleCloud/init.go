package googlecloudstoragecoreclients

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlecloudstoragecoreclientinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/init"
	googlecloudstoragecorewrapperinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/wrappers"
	googlecloudstoragecorewrappers "github.com/horeekaa/backend/core/storages/googleCloud/wrappers"
)

type googleCloudStorageClient struct {
	client       *storage.Client
	pathIdentity string
}

func NewGoogleCloudStorageClient() (googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient, error) {
	return &googleCloudStorageClient{
		pathIdentity: "GoogleCloudStorageClient",
	}, nil
}

func (storageClient *googleCloudStorageClient) Initialize() (bool, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			storageClient.pathIdentity,
			err,
		)
	}
	storageClient.client = client

	return true, nil
}

func (storageClient *googleCloudStorageClient) GetObjectHandle(
	bucketName string,
	objectPath string,
) (googlecloudstoragecorewrapperinterfaces.GCSObjectHandle, error) {
	o := storageClient.client.Bucket(bucketName).Object(objectPath)
	return googlecloudstoragecorewrappers.NewGCSObjectHandle(o)
}

func (storageClient *googleCloudStorageClient) CopyWrite(
	dst io.Writer, src io.Reader,
) (written int64, err error) {
	return io.Copy(dst, src)
}
