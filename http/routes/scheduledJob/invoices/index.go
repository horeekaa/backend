package invoicescheduledjobroutes

import (
	"github.com/go-chi/chi"
)

func Route(r chi.Router) {
	r.Get("/createInvoice", createInvoiceHandler())
}
