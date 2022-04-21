package scheduledjobroutes

import (
	"github.com/go-chi/chi"
	invoicescheduledjobroutes "github.com/horeekaa/backend/http/routes/scheduledJob/invoices"
	purchaseorderscheduledjobroutes "github.com/horeekaa/backend/http/routes/scheduledJob/purchaseOrders"
	purchaseordertosupplyscheduledjobroutes "github.com/horeekaa/backend/http/routes/scheduledJob/purchaseOrdersToSupply"
)

func Route(r chi.Router) {
	r.Get("/createInvoice", invoicescheduledjobroutes.CreateInvoiceHandler())
	r.Get("/updateDueInvoice", invoicescheduledjobroutes.UpdateDueInvoiceHandler())
	r.Get("/processPOToSupply", purchaseordertosupplyscheduledjobroutes.ProcessPOToSupplyHandler())
	r.Get("/updateReceivedPurchaseOrder", purchaseorderscheduledjobroutes.UpdateReceivedPurchaseOrderHandler())
}
