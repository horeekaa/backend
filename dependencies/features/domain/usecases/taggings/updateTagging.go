package taggingpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecases "github.com/horeekaa/backend/features/taggings/domain/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
)

type BulkUpdateTaggingUsecaseDependency struct{}

func (_ *BulkUpdateTaggingUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingRepository,
			approveUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingRepository,
		) taggingpresentationusecaseinterfaces.BulkUpdateTaggingUsecase {
			bulkUpdateTaggingUsecase, _ := taggingpresentationusecases.NewBulkUpdateTaggingUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateTaggingRepo,
				approveUpdateTaggingRepo,
			)
			return bulkUpdateTaggingUsecase
		},
	)
}
