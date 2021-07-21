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

type createDescriptivePhotoTransactionComponent struct {
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	gcsBasicImageStoring       googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
}

func NewCreateDescriptivePhotoTransactionComponent(
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
) (descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent, error) {
	return &createDescriptivePhotoTransactionComponent{
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		gcsBasicImageStoring:       gcsBasicImageStoring,
	}, nil
}

func (createDescPhotoTrx *createDescriptivePhotoTransactionComponent) PreTransaction(
	createDescriptivePhotoInput *model.InternalCreateDescriptivePhoto,
) (*model.InternalCreateDescriptivePhoto, error) {
	return createDescriptivePhotoInput, nil
}

func (createDescPhotoTrx *createDescriptivePhotoTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateDescriptivePhoto,
) (*model.DescriptivePhoto, error) {
	descPhotoToCreate := &model.DatabaseCreateDescriptivePhoto{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, descPhotoToCreate)

	if input.Photo != nil {
		photoUrl, err := createDescPhotoTrx.gcsBasicImageStoring.UploadImage(
			context.Background(),
			input.Category,
			googlecloudstoragecoretypes.GCSFileUpload(*input.Photo),
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createDescriptionPhoto",
				err,
			)
		}
		descPhotoToCreate.PhotoURL = &photoUrl
	}

	createdDescPhoto, err := createDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().Create(
		descPhotoToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createDescriptionPhoto",
			err,
		)
	}

	return createdDescPhoto, nil
}
