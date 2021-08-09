package taggingdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bulkCreateTaggingTransactionComponent struct {
	taggingDataSource             databasetaggingdatasourceinterfaces.TaggingDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	tagDataSource                 databasetagdatasourceinterfaces.TagDataSource
	organizationDataSource        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	structFieldIteratorUtility    coreutilityinterfaces.StructFieldIteratorUtility
	createTaggingUsecaseComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingUsecaseComponent
	generatedObjectID             *primitive.ObjectID
}

func NewBulkCreateTaggingTransactionComponent(
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent, error) {
	return &bulkCreateTaggingTransactionComponent{
		taggingDataSource:          taggingDataSource,
		loggingDataSource:          loggingDataSource,
		tagDataSource:              tagDataSource,
		organizationDataSource:     organizationDataSource,
		productDataSource:          productDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
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
	taggings := []*model.Tagging{}
	taggingsToCreate := []*model.DatabaseCreateTagging{}
	jsonTemp, _ := json.Marshal(input)

	if input.CorrelatedTags != nil {
		for _, correlatedTag := range input.CorrelatedTags {
			taggingToCreate := &model.DatabaseCreateTagging{}
			json.Unmarshal(jsonTemp, taggingToCreate)
			checkedCorrelatedTag, err := bulkCreateTaggingTrx.tagDataSource.GetMongoDataSource().FindByID(
				*correlatedTag.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.CorrelatedTag = &model.ObjectIDOnly{
				ID: &checkedCorrelatedTag.ID,
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
			checkedOrganization, err := bulkCreateTaggingTrx.organizationDataSource.GetMongoDataSource().FindByID(
				*organization.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.Organization = &model.ObjectIDOnly{
				ID: &checkedOrganization.ID,
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
			checkedProduct, err := bulkCreateTaggingTrx.productDataSource.GetMongoDataSource().FindByID(
				*product.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/bulkCreateTagging",
					err,
				)
			}
			taggingToCreate.Product = &model.ObjectIDOnly{
				ID: &checkedProduct.ID,
			}
			taggingToCreate.TaggingType = func(tt model.TaggingType) *model.TaggingType {
				return &tt
			}(model.TaggingTypeProduct)

			taggingsToCreate = append(taggingsToCreate, taggingToCreate)
		}
	}

	for _, taggingToCreate := range taggingsToCreate {
		fieldChanges := []*model.FieldChangeDataInput{}
		bulkCreateTaggingTrx.structFieldIteratorUtility.SetIteratingFunc(
			func(tag interface{}, field interface{}, tagString *interface{}) {
				*tagString = fmt.Sprintf(
					"%v%v",
					*tagString,
					tag,
				)

				fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
					Name:     fmt.Sprint(*tagString),
					Type:     reflect.TypeOf(field).Kind().String(),
					NewValue: fmt.Sprint(field),
				})
				*tagString = ""
			},
		)
		bulkCreateTaggingTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
			func(tag interface{}, tagString *interface{}) {
				*tagString = fmt.Sprintf(
					"%v%v.",
					*tagString,
					tag,
				)
			},
		)
		var tagString interface{} = ""
		bulkCreateTaggingTrx.structFieldIteratorUtility.IterateStruct(
			*taggingToCreate,
			&tagString,
		)

		generatedObjectID := bulkCreateTaggingTrx.GetCurrentObjectID()
		loggingOutput, err := bulkCreateTaggingTrx.loggingDataSource.GetMongoDataSource().Create(
			&model.CreateLogging{
				Collection: "Tagging",
				Document: &model.ObjectIDOnly{
					ID: &generatedObjectID,
				},
				FieldChanges: fieldChanges,
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
