package purchaseorderscheduledjobroutes

import (
	"log"
	"net/http"

	container "github.com/golobby/container/v2"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

func UpdateReceivedPurchaseOrderHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			googleInternalAuth := r.Header.Get("X-Appengine-Cron")
			if googleInternalAuth != "true" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var updatePurchaseOrderUsecase purchaseorderpresentationusecaseinterfaces.UpdatePurchaseOrderUsecase
			container.Make(&updatePurchaseOrderUsecase)

			_, err := updatePurchaseOrderUsecase.Execute(
				purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput{
					CronAuthenticated:   true,
					UpdatePurchaseOrder: &model.UpdatePurchaseOrder{},
				},
			)
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
