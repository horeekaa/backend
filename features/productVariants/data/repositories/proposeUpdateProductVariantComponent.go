package productvariantdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateProductVariantTransactionComponent struct {
	productVariantDataSource               databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	createDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdateProductVariantTransactionComponent(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (productvariantdomainrepositoryinterfaces.ProposeUpdateProductVariantTransactionComponent, error) {
	return &proposeUpdateProductVariantTransactionComponent{
		productVariantDataSource:               productVariantDataSource,
		loggingDataSource:                      loggingDataSource,
		createDescriptivePhotoComponent:        createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent: proposeUpdateDescriptivePhotoComponent,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (updateProdVariantTrx *proposeUpdateProductVariantTransactionComponent) PreTransaction(
	updateProductVariantInput *model.InternalUpdateProductVariant,
) (*model.InternalUpdateProductVariant, error) {
	return updateProductVariantInput, nil
}

func (updateProdVariantTrx *proposeUpdateProductVariantTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateProductVariant *model.InternalUpdateProductVariant,
) (*model.ProductVariant, error) {
	existingProductVariant, err := updateProdVariantTrx.productVariantDataSource.GetMongoDataSource().FindByID(
		*updateProductVariant.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	if updateProductVariant.Photo != nil {
		if updateProductVariant.Photo.ID != nil {
			updateProductVariant.Photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updateProductVariant.ProposalStatus)
			updateProductVariant.Photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updateProductVariant.SubmittingAccount)
			_, err := updateProdVariantTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
				session,
				updateProductVariant.Photo,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateProductVariant",
					err,
				)
			}
		} else {
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(updateProductVariant.Photo)
			json.Unmarshal(jsonTemp, photoToCreate)
			photoToCreate.Category = model.DescriptivePhotoCategoryProductVariant
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingProductVariant.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updateProductVariant.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updateProductVariant.SubmittingAccount)
			if funk.Get(updateProductVariant, "Photo.Photo") != nil {
				photoToCreate.Photo.File = updateProductVariant.Photo.Photo.File
			}
			createdDescriptivePhoto, err := updateProdVariantTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateProductVariant",
					err,
				)
			}

			if existingProductVariant.Photo != nil {
				_, err = updateProdVariantTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					&model.InternalUpdateDescriptivePhoto{
						ID:       &existingProductVariant.Photo.ID,
						IsActive: func(b bool) *bool { return &b }(false),
					},
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updateProductVariant",
						err,
					)
				}
			}
			updateProductVariant.Photo = &model.InternalUpdateDescriptivePhoto{
				ID: &createdDescriptivePhoto.ID,
			}
		}
	}

	newDocumentJson, _ := json.Marshal(*updateProductVariant)
	oldDocumentJson, _ := json.Marshal(*existingProductVariant)
	loggingOutput, err := updateProdVariantTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "ProductVariant",
			Document: &model.ObjectIDOnly{
				ID: &existingProductVariant.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateProductVariant.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateProductVariant.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}
	updateProductVariant.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateProductVariant := &model.DatabaseUpdateProductVariant{
		ID: *updateProductVariant.ID,
	}
	jsonExisting, _ := json.Marshal(existingProductVariant)
	json.Unmarshal(jsonExisting, &fieldsToUpdateProductVariant.ProposedChanges)

	var updateProductVariantMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateProductVariant)
	json.Unmarshal(jsonUpdate, &updateProductVariantMap)

	updateProdVariantTrx.mapProcessorUtility.RemoveNil(updateProductVariantMap)

	jsonUpdate, _ = json.Marshal(updateProductVariantMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateProductVariant.ProposedChanges)

	if updateProductVariant.ProposalStatus != nil {
		fieldsToUpdateProductVariant.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateProductVariant.SubmittingAccount.ID,
		}
		if *updateProductVariant.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateProductVariant)
		}
	}

	updatedProdVariant, err := updateProdVariantTrx.productVariantDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateProductVariant.ID,
		},
		fieldsToUpdateProductVariant,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	return updatedProdVariant, nil
}
