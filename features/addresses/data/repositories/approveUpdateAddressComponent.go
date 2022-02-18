package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
)

type approveUpdateAddressTransactionComponent struct {
	addressDataSource   databaseaddressdatasourceinterfaces.AddressDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdateAddressTransactionComponent(
	AddressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (addressdomainrepositoryinterfaces.ApproveUpdateAddressTransactionComponent, error) {
	return &approveUpdateAddressTransactionComponent{
		addressDataSource:   AddressDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (approveAddrTrx *approveUpdateAddressTransactionComponent) PreTransaction(
	input *model.InternalUpdateAddress,
) (*model.InternalUpdateAddress, error) {
	return input, nil
}

func (approveAddrTrx *approveUpdateAddressTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateAddress,
) (*model.Address, error) {
	updateAddress := &model.DatabaseUpdateAddress{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateAddress)

	existingAddress, err := approveAddrTrx.addressDataSource.GetMongoDataSource().FindByID(
		updateAddress.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateAddress",
			err,
		)
	}
	if existingAddress.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingAddress, nil
	}

	previousLog, err := approveAddrTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingAddress.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateAddress",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateAddress.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateAddress.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveAddrTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateAddress",
			err,
		)
	}

	updateAddress.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateAddress := &model.DatabaseUpdateAddress{
		ID: updateAddress.ID,
	}
	jsonExisting, _ := json.Marshal(existingAddress.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateAddress.ProposedChanges)

	var updateAddressMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateAddress)
	json.Unmarshal(jsonUpdate, &updateAddressMap)

	approveAddrTrx.mapProcessorUtility.RemoveNil(updateAddressMap)

	jsonUpdate, _ = json.Marshal(updateAddressMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateAddress.ProposedChanges)

	if updateAddress.ProposalStatus != nil {
		if *updateAddress.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateAddress.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateAddress)
		}
	}

	updatedAddress, err := approveAddrTrx.addressDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateAddress.ID,
		},
		fieldsToUpdateAddress,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateAddress",
			err,
		)
	}

	return updatedAddress, nil
}
