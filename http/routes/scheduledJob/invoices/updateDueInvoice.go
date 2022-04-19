package invoicescheduledjobroutes

import (
	"log"
	"net/http"

	container "github.com/golobby/container/v2"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

func UpdateDueInvoiceHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			googleInternalAuth := r.Header.Get("X-Appengine-Cron")
			if googleInternalAuth != "true" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var updateInvoiceUsecase invoicepresentationusecaseinterfaces.UpdateInvoiceUsecase
			container.Make(&updateInvoiceUsecase)

			_, err := updateInvoiceUsecase.Execute(
				invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput{
					CronAuthenticated: true,
					UpdateInvoice:     &model.UpdateInvoice{},
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
