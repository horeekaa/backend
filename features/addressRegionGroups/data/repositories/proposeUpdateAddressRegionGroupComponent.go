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

type proposeUpdateAddressRegionGroupTransactionComponent struct {
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
	loggingDataSource            databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility          coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdateAddressRegionGroupTransactionComponent(
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupTransactionComponent, error) {
	return &proposeUpdateAddressRegionGroupTransactionComponent{
		addressRegionGroupDataSource: addressRegionGroupDataSource,
		loggingDataSource:            loggingDataSource,
		mapProcessorUtility:          mapProcessorUtility,
	}, nil
}

func (updateAddressRegionGroupTrx *proposeUpdateAddressRegionGroupTransactionComponent) PreTransaction(
	input *model.InternalUpdateAddressRegionGroup,
) (*model.InternalUpdateAddressRegionGroup, error) {
	return input, nil
}

func (updateAddressRegionGroupTrx *proposeUpdateAddressRegionGroupTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	updateAddressRegionGroup := &model.DatabaseUpdateAddressRegionGroup{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateAddressRegionGroup)

	existingAddressRegionGroup, err := updateAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().FindByID(
		updateAddressRegionGroup.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRegionGroup",
			err,
		)
	}
	newDocumentJson, _ := json.Marshal(*updateAddressRegionGroup)
	oldDocumentJson, _ := json.Marshal(*existingAddressRegionGroup)
	loggingOutput, err := updateAddressRegionGroupTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "AddressRegionGroup",
			Document: &model.ObjectIDOnly{
				ID: &existingAddressRegionGroup.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateAddressRegionGroup.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateAddressRegionGroup.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRegionGroup",
			err,
		)
	}
	updateAddressRegionGroup.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateAddressRegionGroup := &model.DatabaseUpdateAddressRegionGroup{
		ID: updateAddressRegionGroup.ID,
	}
	jsonExisting, _ := json.Marshal(existingAddressRegionGroup)
	json.Unmarshal(jsonExisting, &fieldsToUpdateAddressRegionGroup.ProposedChanges)

	var updateAddressRegionGroupMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateAddressRegionGroup)
	json.Unmarshal(jsonUpdate, &updateAddressRegionGroupMap)

	updateAddressRegionGroupTrx.mapProcessorUtility.RemoveNil(updateAddressRegionGroupMap)

	jsonUpdate, _ = json.Marshal(updateAddressRegionGroupMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateAddressRegionGroup.ProposedChanges)

	if updateAddressRegionGroup.ProposalStatus != nil {
		fieldsToUpdateAddressRegionGroup.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateAddressRegionGroup.SubmittingAccount.ID,
		}
		if *updateAddressRegionGroup.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateAddressRegionGroup)
		}
	}

	updatedAddressRegionGroup, err := updateAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().Update(
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
