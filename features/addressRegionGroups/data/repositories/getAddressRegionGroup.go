package addressregiongroupdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAddressRegionGroupRepository struct {
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
	pathIdentity                 string
}

func NewGetAddressRegionGroupRepository(
	AddressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
) (addressregiongroupdomainrepositoryinterfaces.GetAddressRegionGroupRepository, error) {
	return &getAddressRegionGroupRepository{
		AddressRegionGroupDataSource,
		"GetAddressRegionGroupRepository",
	}, nil
}

func (getAddressRegionGroupRepo *getAddressRegionGroupRepository) Execute(filterFields *model.AddressRegionGroupFilterFields) (*model.AddressRegionGroup, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	addressRegionGroup, err := getAddressRegionGroupRepo.addressRegionGroupDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAddressRegionGroupRepo.pathIdentity,
			err,
		)
	}

	return addressRegionGroup, nil
}
