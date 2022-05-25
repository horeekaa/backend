package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	notificationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/notifications/presentation/usecases"
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) BulkUpdateNotification(ctx context.Context, bulkUpdateNotification model.BulkUpdateNotification) ([]*model.Notification, error) {
	var bulkUpdateNotificationUsecase notificationpresentationusecaseinterfaces.BulkUpdateNotificationUsecase
	container.Make(&bulkUpdateNotificationUsecase)
	return bulkUpdateNotificationUsecase.Execute(
		notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput{
			Context:                ctx,
			BulkUpdateNotification: &bulkUpdateNotification,
		},
	)
}

func (r *notificationResolver) Message(ctx context.Context, obj *model.Notification) (*model.NotificationMessage, error) {
	var CreateNotificationMessage notificationpresentationusecaseinterfaces.CreateNotificationMessage
	container.Make(&CreateNotificationMessage)
	return CreateNotificationMessage.Execute(
		notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput{
			Context:      ctx,
			Notification: obj,
		},
	)
}

func (r *queryResolver) Notifications(ctx context.Context, filterFields model.NotificationFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Notification, error) {
	var getNotificationsUsecase notificationpresentationusecaseinterfaces.GetAllNotificationUsecase
	container.Make(&getNotificationsUsecase)
	return getNotificationsUsecase.Execute(
		notificationpresentationusecasetypes.GetAllNotificationUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Notification returns generated.NotificationResolver implementation.
func (r *Resolver) Notification() generated.NotificationResolver { return &notificationResolver{r} }

type notificationResolver struct{ *Resolver }
