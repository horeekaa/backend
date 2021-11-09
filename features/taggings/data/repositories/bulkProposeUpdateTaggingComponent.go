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
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type bulkProposeUpdateTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility           coreutilityinterfaces.MapProcessorUtility
	updateTaggingUsecaseComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingUsecaseComponent
}

func NewBulkProposeUpdateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent, error) {
	return &bulkProposeUpdateTaggingTransactionComponent{
		taggingDataSource:      taggingDataSource,
		tagDataSource:          tagDataSource,
		organizationDataSource: organizationDataSource,
		productDataSource:      productDataSource,
		loggingDataSource:      loggingDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (bulkProposeUpdateTaggingComp *bulkProposeUpdateTaggingTransactionComponent) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingUsecaseComponent,
) (bool, error) {
	bulkProposeUpdateTaggingComp.updateTaggingUsecaseComponent = usecaseComponent
	return true, nil
}

func (bulkProposeUpdateTaggingComp *bulkProposeUpdateTaggingTransactionComponent) PreTransaction(
	input *model.InternalBulkUpdateTagging,
) (*model.InternalBulkUpdateTagging, error) {
	if bulkProposeUpdateTaggingComp.updateTaggingUsecaseComponent == nil {
		return input, nil
	}
	return bulkProposeUpdateTaggingComp.updateTaggingUsecaseComponent.Validation(input)
}

func (bulkProposeUpdateTaggingComp *bulkProposeUpdateTaggingTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalBulkUpdateTagging,
) ([]*model.Tagging, error) {
	taggings := []*model.Tagging{}
	if input.CorrelatedTag != nil {
		_, err := bulkProposeUpdateTaggingComp.taggingDataSource.GetMongoDataSource().FindByID(
			*input.CorrelatedTag.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}
		input.TaggingType = func(tt model.TaggingType) *model.TaggingType {
			return &tt
		}(model.TaggingTypeTagging)
	}

	if input.Product != nil {
		_, err := bulkProposeUpdateTaggingComp.productDataSource.GetMongoDataSource().FindByID(
			*input.Product.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}
		input.TaggingType = func(tt model.TaggingType) *model.TaggingType {
			return &tt
		}(model.TaggingTypeProduct)
	}

	if input.Organization != nil {
		_, err := bulkProposeUpdateTaggingComp.organizationDataSource.GetMongoDataSource().FindByID(
			*input.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}
		input.TaggingType = func(tt model.TaggingType) *model.TaggingType {
			return &tt
		}(model.TaggingTypeOrganization)
	}

	jsonTemp, _ := json.Marshal(input)
	for _, id := range input.IDs {
		taggingToUpdate := &model.DatabaseUpdateTagging{
			ID: *id,
		}
		json.Unmarshal(jsonTemp, taggingToUpdate)
		existingTagging, err := bulkProposeUpdateTaggingComp.taggingDataSource.GetMongoDataSource().FindByID(
			*id,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}

		newDocumentJson, _ := json.Marshal(*taggingToUpdate)
		oldDocumentJson, _ := json.Marshal(*existingTagging)
		loggingOutput, err := bulkProposeUpdateTaggingComp.loggingDataSource.GetMongoDataSource().Create(
			&model.CreateLogging{
				Collection: "Tagging",
				Document: &model.ObjectIDOnly{
					ID: &existingTagging.ID,
				},
				NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
				OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
				CreatedByAccount: &model.ObjectIDOnly{
					ID: taggingToUpdate.SubmittingAccount.ID,
				},
				Activity:       model.LoggedActivityUpdate,
				ProposalStatus: *taggingToUpdate.ProposalStatus,
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}
		taggingToUpdate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

		fieldsToUpdateTagging := &model.DatabaseUpdateTagging{}
		jsonExisting, _ := json.Marshal(existingTagging)
		json.Unmarshal(jsonExisting, &fieldsToUpdateTagging.ProposedChanges)

		var updateTaggingMap map[string]interface{}
		jsonUpdate, _ := json.Marshal(taggingToUpdate)
		json.Unmarshal(jsonUpdate, &updateTaggingMap)

		bulkProposeUpdateTaggingComp.mapProcessorUtility.RemoveNil(updateTaggingMap)

		jsonUpdate, _ = json.Marshal(updateTaggingMap)
		json.Unmarshal(jsonUpdate, &fieldsToUpdateTagging.ProposedChanges)

		if taggingToUpdate.ProposalStatus != nil {
			fieldsToUpdateTagging.RecentApprovingAccount = &model.ObjectIDOnly{
				ID: taggingToUpdate.SubmittingAccount.ID,
			}
			if *taggingToUpdate.ProposalStatus == model.EntityProposalStatusApproved {
				json.Unmarshal(jsonUpdate, fieldsToUpdateTagging)
			}
		}

		updatedTagging, err := bulkProposeUpdateTaggingComp.taggingDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": fieldsToUpdateTagging.ID,
			},
			fieldsToUpdateTagging,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateTagging",
				err,
			)
		}

		taggings = append(taggings, updatedTagging)
	}

	return taggings, nil
}
