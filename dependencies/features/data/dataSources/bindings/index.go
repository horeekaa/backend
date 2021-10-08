package datasourcesdependencies

import (
	firebaseauthdependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/accounts/authentication"
	mongodbaccountdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/accounts/databases/mongodb"
	mongodbaddressdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/addresses/databases/mongodb"
	mongodbdescriptivephotodatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/descriptivePhotos/databases/mongodb"
	mongodbloggingdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/loggings/databases/mongodb"
	mongodbmemberaccessrefdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/memberAccessRefs/databases/mongodb"
	mongodbmemberaccessdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/memberAccesses/databases/mongodb"
	mongodbmouitemdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/mouItems/databases/mongodb"
	mongodbmoudatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/mous/databases/mongodb"
	mongodbnotificationdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/notifications/databases/mongodb"
	mongodborganizationdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/organizations/databases/mongodb"
	mongodbproductvariantdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/productVariants/databases/mongodb"
	mongodbproductdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/products/databases/mongodb"
	mongodbtaggingdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/taggings/databases/mongodb"
	mongodbtagdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/tags/databases/mongodb"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type DatasourcesDependency struct{}

func (_ *DatasourcesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&firebaseauthdependencies.FirebaseAuthDependency{},
		&mongodbaccountdatasourcedependencies.AccountDataSourceDependency{},
		&mongodbaddressdatasourcedependencies.AddressDataSourceDependency{},
		&mongodbaccountdatasourcedependencies.PersonDataSourceDependency{},
		&mongodbloggingdatasourcedependencies.LoggingDataSourceDependency{},
		&mongodbmemberaccessdatasourcedependencies.MemberAccessDataSourceDependency{},
		&mongodbmemberaccessrefdatasourcedependencies.MemberAccessRefDataSourceDependency{},
		&mongodborganizationdatasourcedependencies.OrganizationDataSourceDependency{},
		&mongodbproductdatasourcedependencies.ProductDataSourceDependency{},
		&mongodbproductvariantdatasourcedependencies.ProductVariantDataSourceDependency{},
		&mongodbdescriptivephotodatasourcedependencies.DescriptivePhotoDataSourceDependency{},
		&mongodbtagdatasourcedependencies.TagDataSourceDependency{},
		&mongodbtaggingdatasourcedependencies.TaggingDataSourceDependency{},
		&mongodbnotificationdatasourcedependencies.NotificationDataSourceDependency{},
		&mongodbmoudatasourcedependencies.MouDataSourceDependency{},
		&mongodbmouitemdatasourcedependencies.MouItemDataSourceDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
