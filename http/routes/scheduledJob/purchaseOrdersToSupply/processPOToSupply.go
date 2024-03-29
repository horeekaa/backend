package purchaseordertosupplyscheduledjobroutes

import (
	"log"
	"net/http"

	container "github.com/golobby/container/v2"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

func ProcessPOToSupplyHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			googleInternalAuth := r.Header.Get("X-Appengine-Cron")
			if googleInternalAuth != "true" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var processPOToSupplyUsecase purchaseordertosupplypresentationusecaseinterfaces.ProcessPurchaseOrderToSupplyUsecase
			container.Make(&processPOToSupplyUsecase)

			_, err := processPOToSupplyUsecase.Execute()
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		},
	)
}
