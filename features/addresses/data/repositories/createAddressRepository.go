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

type createAddressTransactionComponent struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
	addressLoader     addressdomainrepositoryutilityinterfaces.AddressLoader
}

func NewCreateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
) (addressdomainrepositoryinterfaces.CreateAddressTransactionComponent, error) {
	return &createAddressTransactionComponent{
		addressDataSource: addressDataSource,
		addressLoader:     addressLoader,
	}, nil
}

func (createAddrTrx *createAddressTransactionComponent) PreTransaction(
	createaddressInput *model.InternalCreateAddress,
) (*model.InternalCreateAddress, error) {
	return createaddressInput, nil
}

func (createAddrTrx *createAddressTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateAddress,
) (*model.Address, error) {
	addressToCreate := &model.DatabaseCreateAddress{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, addressToCreate)

	_, err := createAddrTrx.addressLoader.Execute(
		session,
		&addressdomainrepositorytypes.LatLngGeocode{
			Latitude:  addressToCreate.Latitude,
			Longitude: addressToCreate.Longitude,
		},
		addressToCreate.ResolvedGeocoding,
		addressToCreate.AddressRegionGroup,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createAddressRepository",
			err,
		)
	}

	createdAddress, err := createAddrTrx.addressDataSource.GetMongoDataSource().Create(
		addressToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createAddressRepository",
			err,
		)
	}

	return createdAddress, nil
}
