package storagecoredependencies

import (
	"github.com/golobby/container/v2"
	coreconfigs "github.com/horeekaa/backend/core/commons/configs"
	googlecloudstoragecoreclients "github.com/horeekaa/backend/core/storages/googleCloud"
	googlecloudstoragecoreclientinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/init"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	googlecloudstoragecoreoperations "github.com/horeekaa/backend/core/storages/googleCloud/operations"
)

type GoogleCloudStorageDependency struct{}

func (_ GoogleCloudStorageDependency) Bind() {
	container.Singleton(
		func() googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient {
			gcsClient, _ := googlecloudstoragecoreclients.NewGoogleCloudStorageClient()
			gcsClient.Initialize()
			return gcsClient
		},
	)

	container.Singleton(
		func(
			gcsClient googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient,
		) googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation {
			gcsImageStoreOps, _ := googlecloudstoragecoreoperations.NewGCSBasicImageStoringOperation(
				gcsClient,
				coreconfigs.GetEnvVariable(coreconfigs.GoogleCloudConfigStorageBucketName),
			)
			return gcsImageStoreOps
		},
	)
}
