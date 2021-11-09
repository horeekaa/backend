package purchaseordertosupplyscheduledjobroutes

import (
	"github.com/go-chi/chi"
)

func Route(r chi.Router) {
	r.Get("/createPOToSupply", createPOToSupplyHandler())
	r.Get("/processPOToSupply", processPOToSupplyHandler())
}
