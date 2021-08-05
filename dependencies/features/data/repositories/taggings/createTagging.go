package taggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositories "github.com/horeekaa/backend/features/taggings/data/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type CreateTaggingDependency struct{}

func (_ *CreateTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
		) taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent {
			createTaggingComponent, _ := taggingdomainrepositories.NewCreateTaggingTransactionComponent(
				taggingDataSource,
				tagDataSource,
				organizationDataSource,
				productDataSource,
			)
			return createTaggingComponent
		},
	)
}
