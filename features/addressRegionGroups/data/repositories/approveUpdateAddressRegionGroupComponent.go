package addressregiongroupdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
)

type approveUpdateAddressRegionGroupTransactionComponent struct {
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
	loggingDataSource            databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility          coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdateAddressRegionGroupTransactionComponent(
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupTransactionComponent, error) {
	return &approveUpdateAddressRegionGroupTransactionComponent{
		addressRegionGroupDataSource: addressRegionGroupDataSource,
		loggingDataSource:            loggingDataSource,
		mapProcessorUtility:          mapProcessorUtility,
	}, nil
}

func (approveAddressRegionGroupTrx *approveUpdateAddressRegionGroupTransactionComponent) PreTransaction(
	input *model.InternalUpdateAddressRegionGroup,
) (*model.InternalUpdateAddressRegionGroup, error) {
	return input, nil
}

func (approveAddressRegionGroupTrx *approveUpdateAddressRegionGroupTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateAddressRegionGroup *model.InternalUpdateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	existingAddressRegionGroup, err := approveAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().FindByID(
		updateAddressRegionGroup.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRegionGroup",
			err,
		)
	}

	previousLog, err := approveAddressRegionGroupTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingAddressRegionGroup.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRegionGroup",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateAddressRegionGroup.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateAddressRegionGroup.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveAddressRegionGroupTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateAddressRegionGroup.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateAddressRegionGroup := &model.DatabaseUpdateAddressRegionGroup{
		ID: updateAddressRegionGroup.ID,
	}
	jsonExisting, _ := json.Marshal(existingAddressRegionGroup.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateAddressRegionGroup.ProposedChanges)

	var updateAddressRegionGroupMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateAddressRegionGroup)
	json.Unmarshal(jsonUpdate, &updateAddressRegionGroupMap)

	approveAddressRegionGroupTrx.mapProcessorUtility.RemoveNil(updateAddressRegionGroupMap)

	jsonUpdate, _ = json.Marshal(updateAddressRegionGroupMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateAddressRegionGroup.ProposedChanges)

	if updateAddressRegionGroup.ProposalStatus != nil {
		if *updateAddressRegionGroup.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateAddressRegionGroup.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateAddressRegionGroup)
		}
	}

	updatedAddressRegionGroup, err := approveAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateAddressRegionGroup.ID,
		},
		fieldsToUpdateAddressRegionGroup,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRegionGroup",
			err,
		)
	}

	return updatedAddressRegionGroup, nil
}
