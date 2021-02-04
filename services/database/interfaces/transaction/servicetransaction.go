package databaseservicetransactioninterfaces

import (
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(session *databaseserviceoperations.ServiceOptions, preOutput interface{}) (interface{}, error)
}

type DatabaseServiceTransaction interface {
	RunTransaction(input interface{}) (interface{}, error)
}
