package notificationdomainrepositoryutilities

import (
	"fmt"
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

	case model.NotificationCategoryInvoiceUpdatedPlain:
		formattedDueDate := input.PayloadOptions.InvoicePayload.Invoice.PaymentDueDate.Format(
			"02/01/2006",
		)
		titleText = localizer.Get(
			"invoices.updated.plain.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"formattedDueDate": formattedDueDate,
			},
		)
		bodyText = localizer.Get(
			"invoices.updated.plain.messages.notification_body",
		)
		break

	case model.NotificationCategoryInvoiceUpdatedPaymentNeeded:
		titleText = localizer.Get(
			"invoices.updated.paymentNeeded.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId": input.PayloadOptions.InvoicePayload.Invoice.PublicID,
			},
		)
		bodyText = localizer.Get(
			"invoices.updated.paymentNeeded.messages.notification_body",
		)
		break

	case model.NotificationCategoryMouCreated:
		titleText = localizer.Get(
			"mous.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId": input.PayloadOptions.MouPayload.Mou.PublicID,
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
				"publicId": input.PayloadOptions.MouPayload.Mou.PublicID,
			},
		)

		bodyText = localizer.Get(
			"mous.updated.messages.notification_body",
		)
		break

	case model.NotificationCategoryMouApproval:
		titleText = localizer.Get(
			"mous.approval.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId":       input.PayloadOptions.MouPayload.Mou.PublicID,
				"proposalStatus": strings.ToLower(input.PayloadOptions.MouPayload.Mou.ProposedChanges.ProposalStatus.String()),
			},
		)

		bodyText = localizer.Get(
			"mous.approval.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderApproval:
		titleText = localizer.Get(
			"purchaseOrders.approval.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId":       input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.PublicID,
				"proposalStatus": strings.ToLower(input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.ProposedChanges.ProposalStatus.String()),
			},
		)

		bodyText = localizer.Get(
			"purchaseOrders.approval.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderCreated:
		flValue := float32(input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.Total)
		fmtValue := fmt.Sprintf("IDR %3.0f", flValue)
		if flValue > 1000.0 {
			flValue = flValue / 1000
			fmtValue = fmt.Sprintf("IDR %3.2fK", flValue)
		}

		titleText = localizer.Get(
			"purchaseOrders.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"value": fmtValue,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrders.created.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderUpdatedReceived:
		titleText = localizer.Get(
			"purchaseOrders.updated.received.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId": input.PayloadOptions.PurchaseOrderPayload.PurchaseOrder.PublicID,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrders.updated.received.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderItemApproval:
		titleText = localizer.Get(
			"purchaseOrderItems.approval.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name":           input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
				"proposalStatus": strings.ToLower(input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProposedChanges.ProposalStatus.String()),
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.approval.messages.notification_body",
		)
		break

	case model.NotificationCategoryPurchaseOrderItemCreated:
		titleText = localizer.Get(
			"purchaseOrderItems.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name":     input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
				"publicId": input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.PurchaseOrder.PublicID,
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.created.messages.notification_body",
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
				"name":   input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.ProductVariant.Product.Name,
				"status": strings.ToLower(input.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.Status.String()),
			},
		)

		bodyText = localizer.Get(
			"purchaseOrderItems.updated.fulfillment.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderApproval:
		titleText = localizer.Get(
			"supplyOrders.approval.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"publicId":       input.PayloadOptions.SupplyOrderPayload.SupplyOrder.PublicID,
				"proposalStatus": strings.ToLower(input.PayloadOptions.SupplyOrderPayload.SupplyOrder.ProposedChanges.ProposalStatus.String()),
			},
		)

		bodyText = localizer.Get(
			"supplyOrders.approval.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderCreated:
		flValue := float32(input.PayloadOptions.SupplyOrderPayload.SupplyOrder.Total)
		fmtValue := fmt.Sprintf("IDR %3.0f", flValue)
		if flValue > 1000.0 {
			flValue = flValue / 1000
			fmtValue = fmt.Sprintf("IDR %3.2fK", flValue)
		}

		titleText = localizer.Get(
			"supplyOrders.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"value": fmtValue,
			},
		)

		bodyText = localizer.Get(
			"supplyOrders.created.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderItemApproval:
		titleText = localizer.Get(
			"supplyOrderItems.approval.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name":           input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.PurchaseOrderToSupply.ProductVariant.Product.Name,
				"proposalStatus": strings.ToLower(input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.ProposedChanges.ProposalStatus.String()),
			},
		)

		bodyText = localizer.Get(
			"supplyOrderItems.approval.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderItemCreated:
		titleText = localizer.Get(
			"supplyOrderItems.created.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name":     input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.PurchaseOrderToSupply.ProductVariant.Product.Name,
				"publicId": input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.SupplyOrder.PublicID,
			},
		)

		bodyText = localizer.Get(
			"supplyOrderItems.created.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderItemPartnerAgreed:
		titleText = localizer.Get(
			"supplyOrderItems.updated.partnerAgreement.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name": input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.PurchaseOrderToSupply.ProductVariant.Product.Name,
			},
		)

		bodyText = localizer.Get(
			"supplyOrderItems.updated.partnerAgreement.messages.notification_body",
		)
		break

	case model.NotificationCategorySupplyOrderItemAccepted:
		titleText = localizer.Get(
			"supplyOrderItems.updated.acceptance.messages.notification_title",
			&golocalizei18ncoretypes.LocalizerReplacement{
				"name":   input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.PurchaseOrderToSupply.ProductVariant.Product.Name,
				"status": strings.ToLower(input.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.Status.String()),
			},
		)

		bodyText = localizer.Get(
			"supplyOrderItems.updated.acceptance.messages.notification_body",
		)
		break

	}

	(*output).Message = &model.NotificationMessage{
		Title: titleText,
		Body:  bodyText,
	}

	return true, nil
}
