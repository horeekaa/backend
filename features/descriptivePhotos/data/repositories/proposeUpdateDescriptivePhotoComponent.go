package descriptivephotodomainrepositories

import (
	"context"
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateDescriptivePhotoTransactionComponent struct {
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	loggingDataSource          databaseloggingdatasourceinterfaces.LoggingDataSource
	gcsBasicImageStoring       googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
	mapProcessorUtility        coreutilityinterfaces.MapProcessorUtility
	pathIdentity               string
}

func NewProposeUpdateDescriptivePhotoTransactionComponent(
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent, error) {
	return &proposeUpdateDescriptivePhotoTransactionComponent{
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		loggingDataSource:          loggingDataSource,
		gcsBasicImageStoring:       gcsBasicImageStoring,
		mapProcessorUtility:        mapProcessorUtility,
		pathIdentity:               "ProposeUpdateDescriptivePhotoComponent",
	}, nil
}

func (updateDescPhotoTrx *proposeUpdateDescriptivePhotoTransactionComponent) PreTransaction(
	UpdateDescriptivePhotoInput *model.InternalUpdateDescriptivePhoto,
) (*model.InternalUpdateDescriptivePhoto, error) {
	return UpdateDescriptivePhotoInput, nil
}

func (updateDescPhotoTrx *proposeUpdateDescriptivePhotoTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateDescriptivePhoto,
) (*model.DescriptivePhoto, error) {
	updateDescriptivePhoto := &model.DatabaseUpdateDescriptivePhoto{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateDescriptivePhoto)

	existingDescriptivePhoto, err := updateDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
		updateDescriptivePhoto.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateDescPhotoTrx.pathIdentity,
			err,
		)
	}

	if input.Photo != nil {
		go func() {
			updateDescPhotoTrx.gcsBasicImageStoring.DeleteImage(
				context.Background(),
				existingDescriptivePhoto.PhotoURL,
			)
		}()
		photoUrl, err := updateDescPhotoTrx.gcsBasicImageStoring.UploadImage(
			context.Background(),
			existingDescriptivePhoto.Category,
			googlecloudstoragecoretypes.GCSFileUpload(*input.Photo),
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updateDescPhotoTrx.pathIdentity,
				err,
			)
		}
		updateDescriptivePhoto.PhotoURL = &photoUrl
	}

	newDocumentJson, _ := json.Marshal(*updateDescriptivePhoto)
	oldDocumentJson, _ := json.Marshal(*existingDescriptivePhoto)
	loggingOutput, err := updateDescPhotoTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "DescriptivePhoto",
			Document: &model.ObjectIDOnly{
				ID: &existingDescriptivePhoto.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateDescriptivePhoto.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateDescriptivePhoto.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateDescPhotoTrx.pathIdentity,
			err,
		)
	}
	updateDescriptivePhoto.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateDescriptivePhoto := &model.DatabaseUpdateDescriptivePhoto{
		ID: updateDescriptivePhoto.ID,
	}
	jsonExisting, _ := json.Marshal(existingDescriptivePhoto)
	json.Unmarshal(jsonExisting, &fieldsToUpdateDescriptivePhoto.ProposedChanges)

	var updateDescriptivePhotoMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateDescriptivePhoto)
	json.Unmarshal(jsonUpdate, &updateDescriptivePhotoMap)

	updateDescPhotoTrx.mapProcessorUtility.RemoveNil(updateDescriptivePhotoMap)

	jsonUpdate, _ = json.Marshal(updateDescriptivePhotoMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateDescriptivePhoto.ProposedChanges)

	if updateDescriptivePhoto.ProposalStatus != nil {
		fieldsToUpdateDescriptivePhoto.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateDescriptivePhoto.SubmittingAccount.ID,
		}
		if *updateDescriptivePhoto.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateDescriptivePhoto)
		}
	}

	updatedDescPhoto, err := updateDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateDescriptivePhoto.ID,
		},
		fieldsToUpdateDescriptivePhoto,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateDescPhotoTrx.pathIdentity,
			err,
		)
	}

	return updatedDescPhoto, nil
}
