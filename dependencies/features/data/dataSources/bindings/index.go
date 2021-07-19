package datasourcesdependencies

import (
	firebaseauthdependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/accounts/authentication"
	mongodbaccountdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/accounts/databases/mongodb"
	mongodbdescriptivephotodatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/descriptivePhotos/databases/mongodb"
	mongodbloggingdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/loggings/databases/mongodb"
	mongodbmemberaccessrefdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/memberAccessRefs/databases/mongodb"
	mongodbmemberaccessdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/memberAccesses/databases/mongodb"
	mongodborganizationdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/organizations/databases/mongodb"
	mongodbproductdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/products/databases/mongodb"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type DatasourcesDependency struct{}

func (_ *DatasourcesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&firebaseauthdependencies.FirebaseAuthDependency{},
		&mongodbaccountdatasourcedependencies.AccountDataSourceDependency{},
		&mongodbaccountdatasourcedependencies.PersonDataSourceDependency{},
		&mongodbloggingdatasourcedependencies.LoggingDataSourceDependency{},
		&mongodbmemberaccessdatasourcedependencies.MemberAccessDataSourceDependency{},
		&mongodbmemberaccessrefdatasourcedependencies.MemberAccessRefDataSourceDependency{},
		&mongodborganizationdatasourcedependencies.OrganizationDataSourceDependency{},
		&mongodbproductdatasourcedependencies.ProductDataSourceDependency{},
		&mongodbdescriptivephotodatasourcedependencies.DescriptivePhotoDataSourceDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
