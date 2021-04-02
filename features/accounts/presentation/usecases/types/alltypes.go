package accountpresentationusecasetypes

import (
	"context"
)

type LoginUsecaseInput struct {
	AuthHeader  string
	DeviceToken string
	Context     context.Context
}
