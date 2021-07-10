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

type createProductTransactionComponent struct {
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	structFieldIteratorUtility    coreutilityinterfaces.StructFieldIteratorUtility
	createProductUsecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent
}

func NewCreateProductTransactionComponent(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (productdomainrepositoryinterfaces.CreateProductTransactionComponent, error) {
	return &createProductTransactionComponent{
		productDataSource:          productDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
	}, nil
}

func (createProductTrx *createProductTransactionComponent) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent,
) (bool, error) {
	createProductTrx.createProductUsecaseComponent = usecaseComponent
	return true, nil
}

func (createProductTrx *createProductTransactionComponent) PreTransaction(
	input *model.InternalCreateProduct,
) (*model.InternalCreateProduct, error) {
	if createProductTrx.createProductUsecaseComponent == nil {
		return input, nil
	}
	return createProductTrx.createProductUsecaseComponent.Validation(input)
}

func (createProductTrx *createProductTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateProduct,
) (*model.Product, error) {
	fieldChanges := []*model.FieldChangeDataInput{}
	createProductTrx.structFieldIteratorUtility.SetIteratingFunc(
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
	createProductTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createProductTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	generatedObjectID := createProductTrx.productDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createProductTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Product",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			FieldChanges: fieldChanges,
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
			"/createProduct",
			err,
		)
	}

	input.ID = generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, &input.ProposedChanges)

	newProduct, err := createProductTrx.productDataSource.GetMongoDataSource().Create(
		input,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createProduct",
			err,
		)
	}
	return newProduct, nil
}