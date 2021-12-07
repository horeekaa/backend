package taggingdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories/utils"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bulkCreateTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	createTaggingUsecaseComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingUsecaseComponent
	taggingLoaderUtility          taggingdomainrepositoryutilityinterfaces.TaggingLoader
	generatedObjectID             *primitive.ObjectID
}

func NewBulkCreateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	taggingLoaderUtility taggingdomainrepositoryutilityinterfaces.TaggingLoader,
) (taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent, error) {
	return &bulkCreateTaggingTransactionComponent{
		taggingDataSource:      taggingDataSource,
		loggingDataSource:      loggingDataSource,
		tagDataSource:          tagDataSource,
		organizationDataSource: organizationDataSource,
		productDataSource:      productDataSource,
		taggingLoaderUtility:   taggingLoaderUtility,
	}, nil
}

func (bulkCreateTaggingTrx *bulkCreateTaggingTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := bulkCreateTaggingTrx.taggingDataSource.GetMongoDataSource().GenerateObjectID()
	bulkCreateTaggingTrx.generatedObjectID = &generatedObjectID
	return *bulkCreateTaggingTrx.generatedObjectID
}

func (bulkCreateTaggingTrx *bulkCreateTaggingTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if bulkCreateTaggingTrx.generatedObjectID == nil {
		generatedObjectID := bulkCreateTaggingTrx.taggingDataSource.GetMongoDataSource().GenerateObjectID()
		bulkCreateTaggingTrx.generatedObjectID = &generatedObjectID
	}
	return *bulkCreateTaggingTrx.generatedObjectID
}

func (bulkCreateTaggingTrx *bulkCreateTaggingTransactionComponent) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingUsecaseComponent,
) (bool, error) {
	bulkCreateTaggingTrx.createTaggingUsecaseComponent = usecaseComponent
	return true, nil
}

func (bulkCreateTaggingTrx *bulkCreateTaggingTransactionComponent) PreTransaction(
	input *model.InternalCreateTagging,
) (*model.InternalCreateTagging, error) {
	if bulkCreateTaggingTrx.createTaggingUsecaseComponent == nil {
		return input, nil
	}
	return bulkCreateTaggingTrx.createTaggingUsecaseComponent.Validation(input)
}

func (bulkCreateTaggingTrx *bulkCreateTaggingTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateTagging,
) ([]*model.Tagging, error) {
	_, err := bulkCreateTaggingTrx.taggingLoaderUtility.TransactionBody(
		session,
		input.Tag,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/bulkCreateTagging",
			err,
		)
	}
	taggings := []*model.Tagging{}
	taggingsToCreate := []*model.DatabaseCreateTagging{}
	jsonTemp, _ := json.Marshal(input)

	if input.CorrelatedTags != nil {
		for _, correlatedTag := range input.CorrelatedTags {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			_, err := bulkCreateTaggingTrx.tagDataSource.GetMongoDataSource().FindByID(
				*correlatedTag.ID,
				session,
			)
			if err != nil && !input.IgnoreTaggedDocumentCheck {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.CorrelatedTag = &model.ObjectIDOnly{
				ID: correlatedTag.ID,
			}
			taggingToCreate.TaggingType = func(tt model.TaggingType) *model.TaggingType {
				return &tt
			}(model.TaggingTypeTagging)

			taggingsToCreate = append(taggingsToCreate, taggingToCreate)
		}
	}

	if input.Organizations != nil {
		for _, organization := range input.Organizations {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			_, err := bulkCreateTaggingTrx.organizationDataSource.GetMongoDataSource().FindByID(
				*organization.ID,
				session,
			)
			if err != nil && !input.IgnoreTaggedDocumentCheck {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.Organization = &model.ObjectIDOnly{
				ID: organization.ID,
			}
			taggingToCreate.TaggingType = func(tt model.TaggingType) *model.TaggingType {
				return &tt
			}(model.TaggingTypeOrganization)

			taggingsToCreate = append(taggingsToCreate, taggingToCreate)
		}
	}

	if input.Products != nil {
		for _, product := range input.Products {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			_, err := bulkCreateTaggingTrx.productDataSource.GetMongoDataSource().FindByID(
				*product.ID,
				session,
			)
			if err != nil && !input.IgnoreTaggedDocumentCheck {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.Product = &model.ObjectIDOnly{
				ID: product.ID,
			}
			taggingToCreate.TaggingType = func(tt model.TaggingType) *model.TaggingType {
				return &tt
			}(model.TaggingTypeProduct)

			taggingsToCreate = append(taggingsToCreate, taggingToCreate)
		}
	}

	for _, taggingToCreate := range taggingsToCreate {
		newDocumentJson, _ := json.Marshal(*taggingToCreate)
		generatedObjectID := bulkCreateTaggingTrx.GetCurrentObjectID()
		loggingOutput, err := bulkCreateTaggingTrx.loggingDataSource.GetMongoDataSource().Create(
			&model.CreateLogging{
				Collection: "Tagging",
				Document: &model.ObjectIDOnly{
					ID: &generatedObjectID,
				},
				NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
				CreatedByAccount: &model.ObjectIDOnly{
					ID: taggingToCreate.SubmittingAccount.ID,
				},
				Activity:       model.LoggedActivityCreate,
				ProposalStatus: *taggingToCreate.ProposalStatus,
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkCreateTagging",
				err,
			)
		}

		taggingToCreate.ID = generatedObjectID
		taggingToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
		if *taggingToCreate.ProposalStatus == model.EntityProposalStatusApproved {
			taggingToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: taggingToCreate.SubmittingAccount.ID}
		}

		jsonTemp, _ := json.Marshal(taggingToCreate)
		json.Unmarshal(jsonTemp, &taggingToCreate.ProposedChanges)

		newTagging, err := bulkCreateTaggingTrx.taggingDataSource.GetMongoDataSource().Create(
			taggingToCreate,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkCreateTagging",
				err,
			)
		}
		taggings = append(taggings, newTagging)
		bulkCreateTaggingTrx.generatedObjectID = nil
	}

	return taggings, nil
}
