package addressregiongroupdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongroupdomainrepositorytypes "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllAddressRegionGroupRepository struct {
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
	mongoQueryBuilder            mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity                 string
}

func NewGetAllAddressRegionGroupRepository(
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (addressregiongroupdomainrepositoryinterfaces.GetAllAddressRegionGroupRepository, error) {
	return &getAllAddressRegionGroupRepository{
		addressRegionGroupDataSource,
		mongoQueryBuilder,
		"GetAllAddressRegionGroupRepository",
	}, nil
}

func (getAllAddressRegionGroupRepo *getAllAddressRegionGroupRepository) Execute(
	input addressregiongroupdomainrepositorytypes.GetAllAddressRegionGroupInput,
) ([]*model.AddressRegionGroup, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllAddressRegionGroupRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	addressRegionGroups, err := getAllAddressRegionGroupRepo.addressRegionGroupDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllAddressRegionGroupRepo.pathIdentity,
			err,
		)
	}

	return addressRegionGroups, nil
}
