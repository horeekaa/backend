package taggingdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type updateTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	updateTaggingUsecaseComponent taggingdomainrepositoryinterfaces.UpdateTaggingUsecaseComponent
}

func NewUpdateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
) (taggingdomainrepositoryinterfaces.UpdateTaggingTransactionComponent, error) {
	return &updateTaggingTransactionComponent{
		taggingDataSource:      taggingDataSource,
		tagDataSource:          tagDataSource,
		organizationDataSource: organizationDataSource,
		productDataSource:      productDataSource,
	}, nil
}

func (updateTaggingTrx *updateTaggingTransactionComponent) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.UpdateTaggingUsecaseComponent,
) (bool, error) {
	updateTaggingTrx.updateTaggingUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateTaggingTrx *updateTaggingTransactionComponent) PreTransaction(
	input *model.InternalBulkUpdateTagging,
) (*model.InternalBulkUpdateTagging, error) {
	if updateTaggingTrx.updateTaggingUsecaseComponent == nil {
		return input, nil
	}
	return updateTaggingTrx.updateTaggingUsecaseComponent.Validation(input)
}

func (updateTaggingTrx *updateTaggingTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalBulkUpdateTagging,
) ([]*model.Tagging, error) {
	taggings := []*model.Tagging{}
	if input.CorrelatedTag != nil {
		_, err := updateTaggingTrx.taggingDataSource.GetMongoDataSource().FindByID(
			*input.CorrelatedTag.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateTagging",
				err,
			)
		}
	}

	if input.Product != nil {
		_, err := updateTaggingTrx.productDataSource.GetMongoDataSource().FindByID(
			*input.Product.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateTagging",
				err,
			)
		}
	}

	if input.Organization != nil {
		_, err := updateTaggingTrx.organizationDataSource.GetMongoDataSource().FindByID(
			*input.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateTagging",
				err,
			)
		}
	}
	for _, id := range input.IDs {
		tagToUpdate := &model.DatabaseUpdateTagging{}
		jsonTemp, _ := json.Marshal(input)
		json.Unmarshal(jsonTemp, tagToUpdate)
		updatedTagging, err := updateTaggingTrx.taggingDataSource.GetMongoDataSource().Update(
			*id,
			tagToUpdate,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateTagging",
				err,
			)
		}

		taggings = append(taggings, updatedTagging)
	}

	return taggings, nil
}
