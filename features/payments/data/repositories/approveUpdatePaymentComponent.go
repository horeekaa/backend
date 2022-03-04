package paymentdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePaymentTransactionComponent struct {
	paymentDataSource   databasepaymentdatasourceinterfaces.PaymentDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdatePaymentTransactionComponent(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent, error) {
	return &approveUpdatePaymentTransactionComponent{
		paymentDataSource:   paymentDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (updatePaymentTrx *approveUpdatePaymentTransactionComponent) PreTransaction(
	input *model.InternalUpdatePayment,
) (*model.InternalUpdatePayment, error) {
	return input, nil
}

func (updatePaymentTrx *approveUpdatePaymentTransactionComponent) TransactionBody(
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
			"/updatePayment",
			err,
		)
	}

	previousLog, err := updatePaymentTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingPayment.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePayment",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updatePayment.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updatePayment.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := updatePaymentTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePayment",
			err,
		)
	}

	updatePayment.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdatePayment := &model.DatabaseUpdatePayment{
		ID: updatePayment.ID,
	}
	jsonExisting, _ := json.Marshal(existingPayment.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePayment.ProposedChanges)

	var updatePaymentMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePayment)
	json.Unmarshal(jsonUpdate, &updatePaymentMap)

	updatePaymentTrx.mapProcessorUtility.RemoveNil(updatePaymentMap)

	jsonUpdate, _ = json.Marshal(updatePaymentMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePayment.ProposedChanges)

	if updatePayment.ProposalStatus != nil {
		if *updatePayment.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdatePayment.ProposedChanges)
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
			"/updatePayment",
			err,
		)
	}

	return updatedPayment, nil
}
