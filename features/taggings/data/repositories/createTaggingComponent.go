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

type createTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	createTaggingUsecaseComponent taggingdomainrepositoryinterfaces.CreateTaggingUsecaseComponent
}

func NewCreateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
) (taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent, error) {
	return &createTaggingTransactionComponent{
		taggingDataSource:      taggingDataSource,
		tagDataSource:          tagDataSource,
		organizationDataSource: organizationDataSource,
		productDataSource:      productDataSource,
	}, nil
}

func (createTaggingTrx *createTaggingTransactionComponent) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.CreateTaggingUsecaseComponent,
) (bool, error) {
	createTaggingTrx.createTaggingUsecaseComponent = usecaseComponent
	return true, nil
}

func (createTaggingTrx *createTaggingTransactionComponent) PreTransaction(
	input *model.InternalCreateTagging,
) (*model.InternalCreateTagging, error) {
	if createTaggingTrx.createTaggingUsecaseComponent == nil {
		return input, nil
	}
	return createTaggingTrx.createTaggingUsecaseComponent.Validation(input)
}

func (createTaggingTrx *createTaggingTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateTagging,
) ([]*model.Tagging, error) {
	taggings := []*model.Tagging{}
	jsonTemp, _ := json.Marshal(input)

	if input.CorrelatedTags != nil {
		for _, correlatedTag := range input.CorrelatedTags {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			checkedCorrelatedTag, err := createTaggingTrx.tagDataSource.GetMongoDataSource().FindByID(
				*correlatedTag.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggingToCreate.CorrelatedTag = &model.ObjectIDOnly{
				ID: &checkedCorrelatedTag.ID,
			}

			newTagging, err := createTaggingTrx.taggingDataSource.GetMongoDataSource().Create(
				taggingToCreate,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggings = append(taggings, newTagging)
		}
	}

	if input.Organizations != nil {
		for _, organization := range input.Organizations {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			checkedOrganization, err := createTaggingTrx.organizationDataSource.GetMongoDataSource().FindByID(
				*organization.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggingToCreate.Organization = &model.ObjectIDOnly{
				ID: &checkedOrganization.ID,
			}

			newTagging, err := createTaggingTrx.taggingDataSource.GetMongoDataSource().Create(
				taggingToCreate,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggings = append(taggings, newTagging)
		}
	}

	if input.Products != nil {
		for _, product := range input.Products {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			checkedProduct, err := createTaggingTrx.productDataSource.GetMongoDataSource().FindByID(
				*product.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggingToCreate.Product = &model.ObjectIDOnly{
				ID: &checkedProduct.ID,
			}

			newTagging, err := createTaggingTrx.taggingDataSource.GetMongoDataSource().Create(
				taggingToCreate,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createTagging",
					err,
				)
			}
			taggings = append(taggings, newTagging)
		}
	}

	return taggings, nil
}
