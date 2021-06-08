package googlecloudstoragecoreclients

import (
	"context"

	"cloud.google.com/go/storage"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlecloudstoragecoreclientinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/init"
)

type googleCloudStorageClient struct {
	client *storage.Client
}

func NewGoogleCloudStorageClient() (googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient, error) {
	return &googleCloudStorageClient{}, nil
}

func (storageClient *googleCloudStorageClient) Initialize() (bool, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"/newGoogleCloudStorage/init",
			err,
		)
	}
	storageClient.client = client
	defer client.Close()

	return true, nil
}

func (storageClient *googleCloudStorageClient) GetClient() (*storage.Client, error) {
	if storageClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newGoogleCloudStorage/GetClient",
			nil,
		)
	}
	return storageClient.client, nil
}
