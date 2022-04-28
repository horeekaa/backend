package paymentdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type paymentLoader struct {
	invoiceDataSource      databaseinvoicedatasourceinterfaces.InvoiceDataSource
	supplyOrderDataSource  databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewPaymentLoader(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (paymentdomainrepositoryutilityinterfaces.PaymentLoader, error) {
	return &paymentLoader{
		invoiceDataSource,
		supplyOrderDataSource,
		organizationDataSource,
	}, nil
}

func (purcOrderLoader *paymentLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	invoice *model.InvoiceForPaymentInput,
	supplyOrder *model.SupplyOrderForPaymentInput,
	organization *model.OrganizationForPaymentInput,
) (bool, error) {
	invoiceLoadedChan := make(chan bool)
	supplyOrderLoadedChan := make(chan bool)
	organizationLoadedChan := make(chan bool)
	errChan := make(chan error)
	go func() {
		if invoice == nil {
			invoiceLoadedChan <- true
			return
		}

		loadedInvoice, err := purcOrderLoader.invoiceDataSource.GetMongoDataSource().FindByID(
			invoice.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedInvoice)
		json.Unmarshal(jsonTemp, invoice)

		invoiceLoadedChan <- true
	}()

	go func() {
		if supplyOrder == nil {
			supplyOrderLoadedChan <- true
			return
		}

		loadedSupplyOrder, err := purcOrderLoader.supplyOrderDataSource.GetMongoDataSource().FindByID(
			supplyOrder.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedSupplyOrder)
		json.Unmarshal(jsonTemp, supplyOrder)

		supplyOrderLoadedChan <- true
	}()

	go func() {
		if organization == nil {
			organizationLoadedChan <- true
			return
		}

		loadedOrganization, err := purcOrderLoader.organizationDataSource.GetMongoDataSource().FindByID(
			organization.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedOrganization)
		json.Unmarshal(jsonTemp, organization)

		organizationLoadedChan <- true
	}()

	for i := 0; i < 3; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-invoiceLoadedChan:
			i++
		case _ = <-organizationLoadedChan:
			i++
		case _ = <-supplyOrderLoadedChan:
			i++
		}
	}

	return true, nil
}
