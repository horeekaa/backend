package databaseaddressregiongroupdatasourceinterfaces

import (
	mongodbaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/mongodb/interfaces"
)

type AddressRegionGroupDataSource interface {
	GetMongoDataSource() mongodbaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSourceMongo) bool
}
