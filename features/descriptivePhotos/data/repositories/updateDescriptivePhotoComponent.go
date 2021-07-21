package descriptivephotodomainrepositories

import (
	"context"
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateDescriptivePhotoTransactionComponent struct {
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	gcsBasicImageStoring       googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
}

func NewUpdateDescriptivePhotoTransactionComponent(
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
) (descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent, error) {
	return &updateDescriptivePhotoTransactionComponent{
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		gcsBasicImageStoring:       gcsBasicImageStoring,
	}, nil
}

func (updateDescPhotoTrx *updateDescriptivePhotoTransactionComponent) PreTransaction(
	UpdateDescriptivePhotoInput *model.InternalUpdateDescriptivePhoto,
) (*model.InternalUpdateDescriptivePhoto, error) {
	return UpdateDescriptivePhotoInput, nil
}

func (updateDescPhotoTrx *updateDescriptivePhotoTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateDescriptivePhoto,
) (*model.DescriptivePhoto, error) {
	existingDescPhoto, err := updateDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
		input.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateDescriptionPhoto",
			err,
		)
	}

	descPhotoToUpdate := &model.DatabaseUpdateDescriptivePhoto{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, descPhotoToUpdate)

	if input.Photo != nil {
		go func() {
			updateDescPhotoTrx.gcsBasicImageStoring.DeleteImage(
				context.Background(),
				existingDescPhoto.PhotoURL,
			)
		}()
		photoUrl, err := updateDescPhotoTrx.gcsBasicImageStoring.UploadImage(
			context.Background(),
			existingDescPhoto.Category,
			googlecloudstoragecoretypes.GCSFileUpload(*input.Photo),
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateDescriptionPhoto",
				err,
			)
		}
		descPhotoToUpdate.PhotoURL = &photoUrl
	}

	updatedDescPhoto, err := updateDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().Update(
		descPhotoToUpdate.ID,
		descPhotoToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateDescriptionPhoto",
			err,
		)
	}

	return updatedDescPhoto, nil
}
