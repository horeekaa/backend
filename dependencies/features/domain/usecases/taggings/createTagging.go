package taggingpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecases "github.com/horeekaa/backend/features/taggings/domain/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
)

type BulkCreateTaggingUsecaseDependency struct{}

func (_ *BulkCreateTaggingUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createTaggingRepo taggingdomainrepositoryinterfaces.BulkCreateTaggingRepository,
		) taggingpresentationusecaseinterfaces.BulkCreateTaggingUsecase {
			bulkTaggingUcase, _ := taggingpresentationusecases.NewBulkCreateTaggingUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createTaggingRepo,
			)
			return bulkTaggingUcase
		},
	)
}
