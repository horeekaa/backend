package paymentdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdatePaymentTransactionComponent struct {
	paymentDataSource   databasepaymentdatasourceinterfaces.PaymentDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	paymentDataLoader   paymentdomainrepositoryutilityinterfaces.PaymentLoader
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	pathIdentity        string
}

func NewProposeUpdatePaymentTransactionComponent(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	paymentDataLoader paymentdomainrepositoryutilityinterfaces.PaymentLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (paymentdomainrepositoryinterfaces.ProposeUpdatePaymentTransactionComponent, error) {
	return &proposeUpdatePaymentTransactionComponent{
		paymentDataSource:   paymentDataSource,
		loggingDataSource:   loggingDataSource,
		paymentDataLoader:   paymentDataLoader,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ProposeUpdatePaymentComponent",
	}, nil
}

func (updatePaymentTrx *proposeUpdatePaymentTransactionComponent) PreTransaction(
	input *model.InternalUpdatePayment,
) (*model.InternalUpdatePayment, error) {
	return input, nil
}

func (updatePaymentTrx *proposeUpdatePaymentTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdatePayment,
) (*model.Payment, error) {
	updatePayment := &model.DatabaseUpdatePayment{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updatePayment)

	existingPayment, err := updatePaymentTrx.paymentDataSource.GetMongoDataSource().FindByID(
		updatePayment.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePaymentTrx.pathIdentity,
			err,
		)
	}

	_, err = updatePaymentTrx.paymentDataLoader.TransactionBody(
		session,
		updatePayment.Invoice,
		updatePayment.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePaymentTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*updatePayment)
	oldDocumentJson, _ := json.Marshal(*existingPayment)
	loggingOutput, err := updatePaymentTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Payment",
			Document: &model.ObjectIDOnly{
				ID: &existingPayment.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updatePayment.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updatePayment.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePaymentTrx.pathIdentity,
			err,
		)
	}
	updatePayment.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	currentTime := time.Now().UTC()
	updatePayment.UpdatedAt = &currentTime

	fieldsToUpdatePayment := &model.DatabaseUpdatePayment{
		ID: updatePayment.ID,
	}
	jsonExisting, _ := json.Marshal(existingPayment)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePayment.ProposedChanges)

	var updatePaymentMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePayment)
	json.Unmarshal(jsonUpdate, &updatePaymentMap)

	updatePaymentTrx.mapProcessorUtility.RemoveNil(updatePaymentMap)

	jsonUpdate, _ = json.Marshal(updatePaymentMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePayment.ProposedChanges)

	if updatePayment.ProposalStatus != nil {
		fieldsToUpdatePayment.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updatePayment.SubmittingAccount.ID,
		}
		if *updatePayment.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatePayment)
		}
	}

	updatedPayment, err := updatePaymentTrx.paymentDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatePayment.ID,
		},
		fieldsToUpdatePayment,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePaymentTrx.pathIdentity,
			err,
		)
	}

	return updatedPayment, nil
}
