package paymentdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createPaymentTransactionComponent struct {
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource
	paymentDataLoader paymentdomainrepositoryutilityinterfaces.PaymentLoader
	generatedObjectID *primitive.ObjectID
	pathIdentity      string
}

func NewCreatePaymentTransactionComponent(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	paymentDataLoader paymentdomainrepositoryutilityinterfaces.PaymentLoader,
) (paymentdomainrepositoryinterfaces.CreatePaymentTransactionComponent, error) {
	return &createPaymentTransactionComponent{
		paymentDataSource: paymentDataSource,
		loggingDataSource: loggingDataSource,
		paymentDataLoader: paymentDataLoader,
		pathIdentity:      "CreatePaymentComponent",
	}, nil
}

func (createPaymentTrx *createPaymentTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createPaymentTrx.paymentDataSource.GetMongoDataSource().GenerateObjectID()
	createPaymentTrx.generatedObjectID = &generatedObjectID
	return *createPaymentTrx.generatedObjectID
}

func (createPaymentTrx *createPaymentTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createPaymentTrx.generatedObjectID == nil {
		generatedObjectID := createPaymentTrx.paymentDataSource.GetMongoDataSource().GenerateObjectID()
		createPaymentTrx.generatedObjectID = &generatedObjectID
	}
	return *createPaymentTrx.generatedObjectID
}

func (createPaymentTrx *createPaymentTransactionComponent) PreTransaction(
	input *model.InternalCreatePayment,
) (*model.InternalCreatePayment, error) {
	return input, nil
}

func (createPaymentTrx *createPaymentTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreatePayment,
) (*model.Payment, error) {
	paymentToCreate := &model.DatabaseCreatePayment{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, paymentToCreate)

	paymentToCreate.Organization = &model.OrganizationForPaymentInput{
		ID: *input.MemberAccess.Organization.ID,
	}

	_, err := createPaymentTrx.paymentDataLoader.TransactionBody(
		session,
		paymentToCreate.Invoice,
		paymentToCreate.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPaymentTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*paymentToCreate)
	generatedObjectID := createPaymentTrx.GetCurrentObjectID()
	loggingOutput, err := createPaymentTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Payment",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: paymentToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *paymentToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPaymentTrx.pathIdentity,
			err,
		)
	}

	paymentToCreate.ID = generatedObjectID
	paymentToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *paymentToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		paymentToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: paymentToCreate.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(paymentToCreate)
	json.Unmarshal(jsonTemp, &paymentToCreate.ProposedChanges)

	newPayment, err := createPaymentTrx.paymentDataSource.GetMongoDataSource().Create(
		paymentToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPaymentTrx.pathIdentity,
			err,
		)
	}
	createPaymentTrx.generatedObjectID = nil

	return newPayment, nil
}
