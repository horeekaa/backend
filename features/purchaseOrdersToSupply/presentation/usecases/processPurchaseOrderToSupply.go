package purchaseordertosupplypresentationusecaseinterfaces

type ProcessPurchaseOrderToSupplyUsecase interface {
	Execute() (bool, error)
}
