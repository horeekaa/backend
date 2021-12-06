package purchaseordertosupplyscheduledjobroutes

import (
	"github.com/go-chi/chi"
)

func Route(r chi.Router) {
	r.Get("/processPOToSupply", processPOToSupplyHandler())
}
