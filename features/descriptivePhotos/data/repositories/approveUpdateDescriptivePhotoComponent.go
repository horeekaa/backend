package descriptivephotodomainrepositories

import (
	"context"
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasedescriptivePhotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
)

type approveUpdateDescriptivePhotoTransactionComponent struct {
	descriptivePhotoDataSource databasedescriptivePhotodatasourceinterfaces.DescriptivePhotoDataSource
	loggingDataSource          databaseloggingdatasourceinterfaces.LoggingDataSource
	gcsBasicImageStoring       googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
	mapProcessorUtility        coreutilityinterfaces.MapProcessorUtility
	pathIdentity               string
}

func NewApproveUpdateDescriptivePhotoTransactionComponent(
	descriptivePhotoDataSource databasedescriptivePhotodatasourceinterfaces.DescriptivePhotoDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent, error) {
	return &approveUpdateDescriptivePhotoTransactionComponent{
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		loggingDataSource:          loggingDataSource,
		gcsBasicImageStoring:       gcsBasicImageStoring,
		mapProcessorUtility:        mapProcessorUtility,
		pathIdentity:               "ApproveUpdateDescPhotoComponent",
	}, nil
}

func (approveDescPhotoTrx *approveUpdateDescriptivePhotoTransactionComponent) PreTransaction(
	input *model.InternalUpdateDescriptivePhoto,
) (*model.InternalUpdateDescriptivePhoto, error) {
	return input, nil
}

func (approveDescPhotoTrx *approveUpdateDescriptivePhotoTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateDescriptivePhoto,
) (*model.DescriptivePhoto, error) {
	updateDescriptivePhoto := &model.DatabaseUpdateDescriptivePhoto{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateDescriptivePhoto)

	existingDescriptivePhoto, err := approveDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
		updateDescriptivePhoto.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveDescPhotoTrx.pathIdentity,
			err,
		)
	}
	if existingDescriptivePhoto.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingDescriptivePhoto, nil
	}

	previousLog, err := approveDescPhotoTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingDescriptivePhoto.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveDescPhotoTrx.pathIdentity,
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateDescriptivePhoto.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateDescriptivePhoto.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveDescPhotoTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveDescPhotoTrx.pathIdentity,
			err,
		)
	}

	updateDescriptivePhoto.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateDescriptivePhoto := &model.DatabaseUpdateDescriptivePhoto{
		ID: updateDescriptivePhoto.ID,
	}
	jsonExisting, _ := json.Marshal(existingDescriptivePhoto.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateDescriptivePhoto.ProposedChanges)

	var updateDescriptivePhotoMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateDescriptivePhoto)
	json.Unmarshal(jsonUpdate, &updateDescriptivePhotoMap)

	approveDescPhotoTrx.mapProcessorUtility.RemoveNil(updateDescriptivePhotoMap)

	jsonUpdate, _ = json.Marshal(updateDescriptivePhotoMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateDescriptivePhoto.ProposedChanges)

	if updateDescriptivePhoto.ProposalStatus != nil {
		if *updateDescriptivePhoto.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateDescriptivePhoto.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateDescriptivePhoto)
		}

		if *updateDescriptivePhoto.ProposalStatus == model.EntityProposalStatusRejected &&
			existingDescriptivePhoto.ProposalStatus == model.EntityProposalStatusProposed {
			go func() {
				approveDescPhotoTrx.gcsBasicImageStoring.DeleteImage(
					context.Background(),
					existingDescriptivePhoto.PhotoURL,
				)
			}()
			fieldsToUpdateDescriptivePhoto.IsActive = func(b bool) *bool { return &b }(false)
		}
	}

	updatedDescriptivePhoto, err := approveDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateDescriptivePhoto.ID,
		},
		fieldsToUpdateDescriptivePhoto,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveDescPhotoTrx.pathIdentity,
			err,
		)
	}

	return updatedDescriptivePhoto, nil
}
