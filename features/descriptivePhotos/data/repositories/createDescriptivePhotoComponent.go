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
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createDescriptivePhotoTransactionComponent struct {
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	loggingDataSource          databaseloggingdatasourceinterfaces.LoggingDataSource
	gcsBasicImageStoring       googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
	generatedObjectID          *primitive.ObjectID
}

func (createDescPhotoTrx *createDescriptivePhotoTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().GenerateObjectID()
	createDescPhotoTrx.generatedObjectID = &generatedObjectID
	return *createDescPhotoTrx.generatedObjectID
}

func (createDescPhotoTrx *createDescriptivePhotoTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createDescPhotoTrx.generatedObjectID == nil {
		generatedObjectID := createDescPhotoTrx.descriptivePhotoDataSource.GetMongoDataSource().GenerateObjectID()
		createDescPhotoTrx.generatedObjectID = &generatedObjectID
	}
	return *createDescPhotoTrx.generatedObjectID
}

func NewCreateDescriptivePhotoTransactionComponent(
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
) (descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent, error) {
	return &createDescriptivePhotoTransactionComponent{
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		loggingDataSource:          loggingDataSource,
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
		input.PhotoURL = &photoUrl
	}

	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createDescPhotoTrx.GetCurrentObjectID()
	loggingOutput, err := createDescPhotoTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "DescriptivePhoto",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: input.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *input.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createDescriptivePhoto",
			err,
		)
	}
	input.ID = &generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	descPhotoToCreate := &model.DatabaseCreateDescriptivePhoto{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, descPhotoToCreate)
	json.Unmarshal(jsonTemp, &descPhotoToCreate.ProposedChanges)

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
	createDescPhotoTrx.generatedObjectID = nil

	return createdDescPhoto, nil
}
