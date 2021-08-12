package taggingpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecases "github.com/horeekaa/backend/features/taggings/domain/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
)

type GetAllTaggingUsecaseDependency struct{}

func (_ GetAllTaggingUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllTaggingRepo taggingdomainrepositoryinterfaces.GetAllTaggingRepository,
		) taggingpresentationusecaseinterfaces.GetAllTaggingUsecase {
			getAllTaggingUcase, _ := taggingpresentationusecases.NewGetAllTaggingUsecase(
				getAccountFromAuthDataRepo,
				getAccountRepo,
				getAllTaggingRepo,
			)
			return getAllTaggingUcase
		},
	)
}
