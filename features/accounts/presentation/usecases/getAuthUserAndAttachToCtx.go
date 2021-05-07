package accountpresentationusecaseinterfaces

import (
	"context"

	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
)

type GetAuthUserAndAttachToCtxUsecase interface {
	Execute(input accountpresentationusecasetypes.GetAuthUserAndAttachToCtxInput) (context.Context, error)
}
