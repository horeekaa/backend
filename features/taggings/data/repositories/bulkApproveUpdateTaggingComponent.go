package taggingdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories/utils"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type bulkApproveUpdateTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility           coreutilityinterfaces.MapProcessorUtility
	taggingLoaderUtility          taggingdomainrepositoryutilityinterfaces.TaggingLoader
	updateTaggingUsecaseComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingUsecaseComponent
}

func NewBulkApproveUpdateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	taggingLoaderUtility taggingdomainrepositoryutilityinterfaces.TaggingLoader,
) (taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent, error) {
	return &bulkApproveUpdateTaggingTransactionComponent{
		taggingDataSource:      taggingDataSource,
		tagDataSource:          tagDataSource,
		organizationDataSource: organizationDataSource,
		productDataSource:      productDataSource,
		loggingDataSource:      loggingDataSource,
		mapProcessorUtility:    mapProcessorUtility,
		taggingLoaderUtility:   taggingLoaderUtility,
	}, nil
}

func (bulkApproveUpdateTaggingComp *bulkApproveUpdateTaggingTransactionComponent) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingUsecaseComponent,
) (bool, error) {
	bulkApproveUpdateTaggingComp.updateTaggingUsecaseComponent = usecaseComponent
	return true, nil
}

func (bulkApproveUpdateTaggingComp *bulkApproveUpdateTaggingTransactionComponent) PreTransaction(
	input *model.InternalBulkUpdateTagging,
) (*model.InternalBulkUpdateTagging, error) {
	if bulkApproveUpdateTaggingComp.updateTaggingUsecaseComponent == nil {
		return input, nil
	}
	return bulkApproveUpdateTaggingComp.updateTaggingUsecaseComponent.Validation(input)
}

func (bulkApproveUpdateTaggingComp *bulkApproveUpdateTaggingTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalBulkUpdateTagging,
) ([]*model.Tagging, error) {
	_, err := bulkApproveUpdateTaggingComp.taggingLoaderUtility.TransactionBody(
		session,
		input.Tag,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/bulkApproveUpdateTagging",
			err,
		)
	}
	taggings := []*model.Tagging{}
	if input.CorrelatedTag != nil {
		_, err := bulkApproveUpdateTaggingComp.taggingDataSource.GetMongoDataSource().FindByID(
			*input.CorrelatedTag.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}
	}

	if input.Product != nil {
		_, err := bulkApproveUpdateTaggingComp.productDataSource.GetMongoDataSource().FindByID(
			*input.Product.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}
	}

	if input.Organization != nil {
		_, err := bulkApproveUpdateTaggingComp.organizationDataSource.GetMongoDataSource().FindByID(
			*input.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}
	}

	jsonTemp, _ := json.Marshal(input)
	for _, id := range input.IDs {
		taggingToUpdate := &model.DatabaseUpdateTagging{
			ID: *id,
		}
		json.Unmarshal(jsonTemp, taggingToUpdate)
		existingTagging, err := bulkApproveUpdateTaggingComp.taggingDataSource.GetMongoDataSource().FindByID(
			*id,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}
		if existingTagging.ProposalStatus == model.EntityProposalStatusApproved {
			taggings = append(taggings, existingTagging)
			continue
		}

		previousLog, err := bulkApproveUpdateTaggingComp.loggingDataSource.GetMongoDataSource().FindByID(
			existingTagging.RecentLog.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}

		logToCreate := &model.CreateLogging{
			Collection: previousLog.Collection,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: taggingToUpdate.RecentApprovingAccount.ID,
			},
			Activity:       previousLog.Activity,
			ProposalStatus: *taggingToUpdate.ProposalStatus,
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"NewDocumentJSON": previousLog.NewDocumentJSON,
				"OldDocumentJSON": previousLog.OldDocumentJSON,
			},
		)
		json.Unmarshal(jsonTemp, logToCreate)

		createdLog, err := bulkApproveUpdateTaggingComp.loggingDataSource.GetMongoDataSource().Create(
			logToCreate,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}

		taggingToUpdate.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

		fieldsToUpdateTagging := &model.DatabaseUpdateTagging{}
		jsonExisting, _ := json.Marshal(existingTagging.ProposedChanges)
		json.Unmarshal(jsonExisting, &fieldsToUpdateTagging.ProposedChanges)

		var updateTaggingMap map[string]interface{}
		jsonUpdate, _ := json.Marshal(taggingToUpdate)
		json.Unmarshal(jsonUpdate, &updateTaggingMap)

		bulkApproveUpdateTaggingComp.mapProcessorUtility.RemoveNil(updateTaggingMap)

		jsonUpdate, _ = json.Marshal(updateTaggingMap)
		json.Unmarshal(jsonUpdate, &fieldsToUpdateTagging.ProposedChanges)

		if taggingToUpdate.ProposalStatus != nil {
			if *taggingToUpdate.ProposalStatus == model.EntityProposalStatusApproved {
				jsonUpdate, _ := json.Marshal(fieldsToUpdateTagging.ProposedChanges)
				json.Unmarshal(jsonUpdate, fieldsToUpdateTagging)
			}
		}

		updatedTagging, err := bulkApproveUpdateTaggingComp.taggingDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": fieldsToUpdateTagging.ID,
			},
			fieldsToUpdateTagging,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkApproveUpdateTagging",
				err,
			)
		}

		taggings = append(taggings, updatedTagging)
	}

	return taggings, nil
}
