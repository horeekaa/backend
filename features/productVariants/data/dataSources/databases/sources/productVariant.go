package databaseproductvariantdatasources

import (
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	mongodbproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/mongodb/interfaces"
)

type productVariantDataSource struct {
	productVariantDataSourceRepoMongo mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo
}

func (descPhotoDataSource *productVariantDataSource) SetMongoDataSource(mongoDataSource mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo) bool {
	descPhotoDataSource.productVariantDataSourceRepoMongo = mongoDataSource
	return true
}

func (descPhotoDataSource *productVariantDataSource) GetMongoDataSource() mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo {
	return descPhotoDataSource.productVariantDataSourceRepoMongo
}

func NewProductVariantDataSource() (databaseproductvariantdatasourceinterfaces.ProductVariantDataSource, error) {
	return &productVariantDataSource{}, nil
}
