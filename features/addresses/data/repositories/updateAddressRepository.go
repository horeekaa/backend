package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type updateAddressTransactionComponent struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
	addressLoader     addressdomainrepositoryutilityinterfaces.AddressLoader
}

func NewUpdateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
) (addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent, error) {
	return &updateAddressTransactionComponent{
		addressDataSource: addressDataSource,
		addressLoader:     addressLoader,
	}, nil
}

func (updateAddrTrx *updateAddressTransactionComponent) PreTransaction(
	updateAddressInput *model.InternalUpdateAddress,
) (*model.InternalUpdateAddress, error) {
	return updateAddressInput, nil
}

func (updateAddrTrx *updateAddressTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateAddress,
) (*model.Address, error) {
	_, err := updateAddrTrx.addressDataSource.GetMongoDataSource().FindByID(
		*input.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddress",
			err,
		)
	}

	addressToUpdate := &model.DatabaseUpdateAddress{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, addressToUpdate)

	if addressToUpdate.Latitude != nil && addressToUpdate.Longitude != nil {
		_, err := updateAddrTrx.addressLoader.Execute(
			&addressdomainrepositorytypes.LatLngGeocode{
				Latitude:  *addressToUpdate.Latitude,
				Longitude: *addressToUpdate.Longitude,
			},
			addressToUpdate.ResolvedGeocoding,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateAddressRepository",
				err,
			)
		}
	}

	updatedAddress, err := updateAddrTrx.addressDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": addressToUpdate.ID,
		},
		addressToUpdate,
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
