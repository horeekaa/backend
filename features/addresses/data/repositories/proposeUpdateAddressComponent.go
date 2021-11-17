package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateAddressTransactionComponent struct {
	addressDataSource   databaseaddressdatasourceinterfaces.AddressDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	addressLoader       addressdomainrepositoryutilityinterfaces.AddressLoader
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (addressdomainrepositoryinterfaces.ProposeUpdateAddressTransactionComponent, error) {
	return &proposeUpdateAddressTransactionComponent{
		addressDataSource:   addressDataSource,
		loggingDataSource:   loggingDataSource,
		addressLoader:       addressLoader,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (updateAddrTrx *proposeUpdateAddressTransactionComponent) PreTransaction(
	updateAddressInput *model.InternalUpdateAddress,
) (*model.InternalUpdateAddress, error) {
	return updateAddressInput, nil
}

func (updateAddrTrx *proposeUpdateAddressTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateAddress *model.InternalUpdateAddress,
) (*model.Address, error) {
	existingAddress, err := updateAddrTrx.addressDataSource.GetMongoDataSource().FindByID(
		*updateAddress.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddress",
			err,
		)
	}

	if updateAddress.Latitude != nil && updateAddress.Longitude != nil {
		updateAddress.ResolvedGeocoding = &model.ResolvedGeocodingInput{}
		updateAddress.AddressRegionGroup = &model.AddressRegionGroupForAddressInput{}
		_, err := updateAddrTrx.addressLoader.Execute(
			session,
			&addressdomainrepositorytypes.LatLngGeocode{
				Latitude:  *updateAddress.Latitude,
				Longitude: *updateAddress.Longitude,
			},
			updateAddress.ResolvedGeocoding,
			updateAddress.AddressRegionGroup,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/proposeUpdateAddressComponent",
				err,
			)
		}
	}

	newDocumentJson, _ := json.Marshal(*updateAddress)
	oldDocumentJson, _ := json.Marshal(*existingAddress)
	loggingOutput, err := updateAddrTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Address",
			Document: &model.ObjectIDOnly{
				ID: &existingAddress.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateAddress.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateAddress.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateAddressComponent",
			err,
		)
	}
	updateAddress.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateAddress := &model.DatabaseUpdateAddress{
		ID: *updateAddress.ID,
	}
	jsonExisting, _ := json.Marshal(existingAddress)
	json.Unmarshal(jsonExisting, &fieldsToUpdateAddress.ProposedChanges)

	var updateAddressMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateAddress)
	json.Unmarshal(jsonUpdate, &updateAddressMap)

	updateAddrTrx.mapProcessorUtility.RemoveNil(updateAddressMap)

	jsonUpdate, _ = json.Marshal(updateAddressMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateAddress.ProposedChanges)

	if updateAddress.ProposalStatus != nil {
		fieldsToUpdateAddress.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateAddress.SubmittingAccount.ID,
		}
		if *updateAddress.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateAddress)
		}
	}

	updatedAddress, err := updateAddrTrx.addressDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateAddress.ID,
		},
		fieldsToUpdateAddress,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddressRepository",
			err,
		)
	}

	return updatedAddress, nil
}
