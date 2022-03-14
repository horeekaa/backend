package googlecloudstoragecoreoperations

import (
	"context"
	"fmt"
	"strings"
	"time"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlecloudstoragecoreclientinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/init"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
	"github.com/horeekaa/backend/model"
	uuid "github.com/satori/go.uuid"
)

type gcsBasicImageStoringOperation struct {
	gcsClient    googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient
	bucketName   string
	pathIdentity string
}

func NewGCSBasicImageStoringOperation(
	gcsClient googlecloudstoragecoreclientinterfaces.GoogleCloudStorageClient,
	bucketName string,
) (googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation, error) {
	return &gcsBasicImageStoringOperation{
		gcsClient,
		bucketName,
		"GoogleCloudStorageBasicOperation",
	}, nil
}

func (gcsBscImageStoringOps *gcsBasicImageStoringOperation) UploadImage(
	ctx context.Context,
	category model.DescriptivePhotoCategory,
	file googlecloudstoragecoretypes.GCSFileUpload,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	objectPath := fmt.Sprintf("images/%s/%s.jpg", category, uuid.NewV4())

	o, _ := gcsBscImageStoringOps.gcsClient.GetObjectHandle(
		gcsBscImageStoringOps.bucketName,
		objectPath,
	)

	wc := o.NewWriter(ctx)
	if _, err := gcsBscImageStoringOps.gcsClient.CopyWrite(wc, file.File); err != nil {
		return "", horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.StoringImageFailed,
			gcsBscImageStoringOps.pathIdentity,
			err,
		)
	}

	if err := wc.Close(); err != nil {
		return "", horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClosingImageStoringWriterFailed,
			gcsBscImageStoringOps.pathIdentity,
			err,
		)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcsBscImageStoringOps.bucketName, objectPath),
		nil
}

func (gcsBscImageStoringOps *gcsBasicImageStoringOperation) DeleteImage(
	ctx context.Context,
	url string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	splittedURL := strings.Split(url, "/")
	objectPath := fmt.Sprintf(
		"images/%s/%s",
		splittedURL[len(splittedURL)-2],
		splittedURL[len(splittedURL)-1],
	)

	o, _ := gcsBscImageStoringOps.gcsClient.GetObjectHandle(
		gcsBscImageStoringOps.bucketName,
		objectPath,
	)
	if err := o.Delete(ctx); err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DeleteImageFailed,
			gcsBscImageStoringOps.pathIdentity,
			err,
		)
	}
	return true, nil
}
