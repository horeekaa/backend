package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateAddressTransactionComponent struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewUpdateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent, error) {
	return &updateAddressTransactionComponent{
		addressDataSource: addressDataSource,
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

	updatedAddress, err := updateAddrTrx.addressDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": addressToUpdate.ID,
		},
		addressToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateAddress",
			err,
		)
	}

	return updatedAddress, nil
}
