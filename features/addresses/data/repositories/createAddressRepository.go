package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createAddressTransactionComponent struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewCreateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (addressdomainrepositoryinterfaces.CreateAddressTransactionComponent, error) {
	return &createAddressTransactionComponent{
		addressDataSource: addressDataSource,
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
