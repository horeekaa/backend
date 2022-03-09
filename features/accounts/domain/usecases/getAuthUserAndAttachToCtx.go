package accountpresentationusecases

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
)

type getAuthUserAndAttachToCtxUsecase struct {
	getUserFromAuthHeaderRepo accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository
	pathIdentity              string
}

func NewGetAuthUserAndAttachToCtxUsecase(
	getUserFromAuthHeaderRepo accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository,
) (accountpresentationusecaseinterfaces.GetAuthUserAndAttachToCtxUsecase, error) {
	return &getAuthUserAndAttachToCtxUsecase{
		getUserFromAuthHeaderRepo,
		"GetAuthUserAndAttachToctx",
	}, nil
}

func (getAuthUserAndAttachToCtx *getAuthUserAndAttachToCtxUsecase) Execute(
	input accountpresentationusecasetypes.GetAuthUserAndAttachToCtxInput,
) (context.Context, error) {
	user, err := getAuthUserAndAttachToCtx.getUserFromAuthHeaderRepo.Execute(
		accountdomainrepositorytypes.GetUserFromAuthHeaderInput{
			AuthHeader: input.AuthHeader,
			Context:    input.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAuthUserAndAttachToCtx.pathIdentity,
			err,
		)
	}

	ctx := context.WithValue(
		input.Context,
		authenticationcoremodels.UserContextKey,
		user,
	)
	return ctx, nil
}
