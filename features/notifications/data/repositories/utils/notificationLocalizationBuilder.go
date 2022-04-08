package notificationdomainrepositoryutilities

import (
	"strings"

	golocalizei18ncoreclientinterfaces "github.com/horeekaa/backend/core/i18n/go-localize/interfaces/init"
	golocalizei18ncoretypes "github.com/horeekaa/backend/core/i18n/go-localize/types"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type notificationLocalizationBuilder struct {
	goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient
}

func NewNotificationLocalizationBuilder(
	goLocalizeI18N golocalizei18ncoreclientinterfaces.GoLocalizeI18NClient,
) (notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder, error) {
	return &notificationLocalizationBuilder{
		goLocalizeI18N: goLocalizeI18N,
	}, nil
}

func (notifLocalBuilder *notificationLocalizationBuilder) Execute(
	input *model.DatabaseNotification,
	output *model.Notification,
	language string,
) (bool, error) {
	notifLocalBuilder.goLocalizeI18N.Initialize(
		strings.ToLower(language),
		"id",
	)
	localizer, _ := notifLocalBuilder.goLocalizeI18N.GetLocalizer()

	titleText := ""
	bodyText := ""
	switch input.NotificationCategory {
	case model.NotificationCategoryMemberAccessInvitationAccepted:
		titleText = localizer.Get(
			"memberAccesses.invitationAccepted.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"personName": input.PayloadOptions.MemberAccessInvitationPayload.MemberAccess.Account.Person.FirstName,
			},
		)
		bodyText = localizer.Get(
			"memberAccesses.invitationAccepted.messages.notification_body",
		)
		break

	case model.NotificationCategoryMemberAccessInvitationRequest:
		titleText = localizer.Get(
			"memberAccesses.invitationRequest.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"submitterName": input.PayloadOptions.MemberAccessInvitationPayload.MemberAccess.SubmittingAccount.Person.FirstName,
				"orgName":       input.PayloadOptions.MemberAccessInvitationPayload.MemberAccess.Organization.Name,
			},
		)
		bodyText = localizer.Get(
			"memberAccesses.invitationRequest.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderSupplyBroadcast:
		titleText = localizer.Get(
			"purchaseOrdersToSupply.orderBroadcast.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"tagName": input.PayloadOptions.PurchaseOrderToSupplyBroadcastPayload.BroadcastedByTag.Name,
			},
		)
		bodyText = localizer.Get(
			"purchaseOrdersToSupply.orderBroadcast.messages.notification_body",
		)
		break
	case model.NotificationCategoryInvoiceCreated:
		formattedDueDate := input.PayloadOptions.InvoicePayload.Invoice.PaymentDueDate.Format(
			"02/01/2006",
		)
		titleText = localizer.Get(
			"invoices.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"formattedDueDate": formattedDueDate,
			},
		)
		bodyText = localizer.Get(
			"invoices.created.messages.notification_body",
		)
		break

	case model.NotificationCategoryInvoiceUpdated:
		formattedDueDate := input.PayloadOptions.InvoicePayload.Invoice.PaymentDueDate.Format(
			"02/01/2006",
		)
		titleText = localizer.Get(
			"invoices.updated.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"formattedDueDate": formattedDueDate,
			},
		)
		bodyText = localizer.Get(
			"invoices.updated.messages.notification_body",
		)
		break

	case model.NotificationCategoryMouCreated:
		titleText = localizer.Get(
			"mous.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"mouPublicId": *input.PayloadOptions.MouPayload.Mou.PublicID,
			},
		)

		bodyText = localizer.Get(
			"mous.created.messages.notification_body",
		)
		break

	case model.NotificationCategoryMouUpdated:
		titleText = localizer.Get(
			"mous.updated.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"mouPublicId": *input.PayloadOptions.MouPayload.Mou.PublicID,
			},
		)

		bodyText = localizer.Get(
			"mous.updated.messages.notification_body",
		)
		break

	case model.NotificationCategoryMouApproved:
		titleText = localizer.Get(
			"mous.approved.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"mouPublicId": *input.PayloadOptions.MouPayload.Mou.PublicID,
			},
		)

		bodyText = localizer.Get(
			"mous.approved.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderApproved:
		titleText = localizer.Get(
			"purchaseOrders.approved.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId": input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.PublicID,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrders.approved.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderCreated:
		titleText = localizer.Get(
			"purchaseOrders.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"value": input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.Total,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrders.created.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderItemApproved:
		titleText = localizer.Get(
			"purchaseOrderItems.approved.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name": input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.approved.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderItemCustomerAgreed:
		titleText = localizer.Get(
			"purchaseOrderItems.updated.customerAgreement.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name": input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.updated.customerAgreement.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderItemFulfilled:
		titleText = localizer.Get(
			"purchaseOrderItems.updated.fulfillment.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name": input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.updated.fulfillment.messages.notification_body",
		)
		break

	}

	(*output).Message = &model.NotificationMessage{
		Title: titleText,
		Body:  bodyText,
	}

	return true, nil
}
