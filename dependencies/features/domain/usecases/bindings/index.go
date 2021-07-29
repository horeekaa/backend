package usecasesdependencies

import (
	accountpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/accounts"
	descriptivephotopresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/descriptivePhotos"
	loggingpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/loggings"
	memberaccessrefpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccessRefs"
	memberaccesspresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccesses"
	organizationpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/organizations"
	productvariantpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/productVariants"
	productpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/products"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type UsecasesDependency struct{}

func (_ *UsecasesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&accountpresentationusecasedependencies.GetAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.GetPersonDataFromAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.LoginUsecaseDependency{},
		&accountpresentationusecasedependencies.LogoutUsecaseDependency{},
		&accountpresentationusecasedependencies.GetAuthUserAndAttachToCtxUsecaseDependency{},

		&loggingpresentationusecasedependencies.GetLoggingUsecaseDependency{},

		&memberaccesspresentationusecasedependencies.CreateMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetAllMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.UpdateMemberAccessUsecaseDependency{},

		&memberaccessrefpresentationusecasedependencies.CreateMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetAllMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.UpdateMemberAccessRefUsecaseDependency{},

		&organizationpresentationusecasedependencies.CreateOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetAllOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.UpdateOrganizationUsecaseDependency{},

		&descriptivephotopresentationusecasedependencies.GetDescriptivePhotoUsecaseDependency{},

		&productpresentationusecasedependencies.CreateProductUsecaseDependency{},
		&productpresentationusecasedependencies.GetAllProductUsecaseDependency{},
		&productpresentationusecasedependencies.GetProductUsecaseDependency{},
		&productpresentationusecasedependencies.UpdateProductUsecaseDependency{},

		&productvariantpresentationusecasedependencies.GetProductVariantUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
