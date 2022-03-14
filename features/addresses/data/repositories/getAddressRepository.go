package addressdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAddressRepository struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
	pathIdentity      string
}

func NewGetAddressRepository(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (addressdomainrepositoryinterfaces.GetAddressRepository, error) {
	return &getAddressRepository{
		addressDataSource,
		"GetAddressRepository",
	}, nil
}

func (getAddressRefRepo *getAddressRepository) Execute(filterFields *model.AddressFilterFields) (*model.Address, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	address, err := getAddressRefRepo.addressDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAddressRefRepo.pathIdentity,
			err,
		)
	}

	return address, nil
}
