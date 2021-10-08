package databaseaddressdatasources

import (
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	mongodbaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/mongodb/interfaces"
)

type addressDataSource struct {
	addressDataSourceMongo mongodbaddressdatasourceinterfaces.AddressDataSourceMongo
}

func (addrDataSource *addressDataSource) SetMongoDataSource(mongoDataSource mongodbaddressdatasourceinterfaces.AddressDataSourceMongo) bool {
	addrDataSource.addressDataSourceMongo = mongoDataSource
	return true
}

func (addrDataSource *addressDataSource) GetMongoDataSource() mongodbaddressdatasourceinterfaces.AddressDataSourceMongo {
	return addrDataSource.addressDataSourceMongo
}

func NewAddressDataSource() (databaseaddressdatasourceinterfaces.AddressDataSource, error) {
	return &addressDataSource{}, nil
}
