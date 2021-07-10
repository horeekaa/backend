package productdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateProductTransactionComponent struct {
	productDataSource                    databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource                    databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                  coreutilityinterfaces.MapProcessorUtility
	structComparisonUtility              coreutilityinterfaces.StructComparisonUtility
	proposeUpdateProductUsecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent
}

func NewProposeUpdateProductTransactionComponent(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
) (productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent, error) {
	return &proposeUpdateProductTransactionComponent{
		productDataSource:       productDataSource,
		loggingDataSource:       loggingDataSource,
		mapProcessorUtility:     mapProcessorUtility,
		structComparisonUtility: structComparisonUtility,
	}, nil
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent,
) (bool, error) {
	updateProdTrx.proposeUpdateProductUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) PreTransaction(
	input *model.InternalUpdateProduct,
) (*model.InternalUpdateProduct, error) {
	if updateProdTrx.proposeUpdateProductUsecaseComponent == nil {
		return input, nil
	}
	return updateProdTrx.proposeUpdateProductUsecaseComponent.Validation(input)
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateProduct *model.InternalUpdateProduct,
) (*model.Product, error) {
	existingProduct, err := updateProdTrx.productDataSource.GetMongoDataSource().FindByID(
		updateProduct.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}
	fieldChanges := []*model.FieldChangeDataInput{}

	updateProdTrx.structComparisonUtility.SetComparisonFunc(
		func(tag interface{}, field1 interface{}, field2 interface{}, tagString *interface{}) {
			if field1 == field2 {
				return
			}
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field1).Kind().String(),
				OldValue: fmt.Sprint(field2),
				NewValue: fmt.Sprint(field1),
			})
			*tagString = ""
		},
	)
	updateProdTrx.structComparisonUtility.SetPreDeepComparisonFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	updateProdTrx.structComparisonUtility.CompareStructs(
		*updateProduct,
		*existingProduct,
		&tagString,
	)

	loggingOutput, err := updateProdTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Product",
			Document: &model.ObjectIDOnly{
				ID: &existingProduct.ID,
			},
			FieldChanges: fieldChanges,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateProduct.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateProduct.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}
	updateProduct.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateProduct := &model.InternalUpdateProduct{
		ID: updateProduct.ID,
	}
	jsonExisting, _ := json.Marshal(existingProduct)
	json.Unmarshal(jsonExisting, &fieldsToUpdateProduct.ProposedChanges)

	var updateProductMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateProduct)
	json.Unmarshal(jsonUpdate, &updateProductMap)

	updateProdTrx.mapProcessorUtility.RemoveNil(updateProductMap)

	jsonUpdate, _ = json.Marshal(updateProductMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateProduct.ProposedChanges)

	if updateProduct.ProposalStatus != nil {
		if *updateProduct.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateProduct.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateProduct)
		}
	}

	updatedProduct, err := updateProdTrx.productDataSource.GetMongoDataSource().Update(
		fieldsToUpdateProduct.ID,
		fieldsToUpdateProduct,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}

	return updatedProduct, nil
}
