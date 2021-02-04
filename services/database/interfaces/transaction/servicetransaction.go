package databaseservicetransactioninterface

import (
	databaseserviceinterface "github.com/horeekaa/backend/services/database/interfaces/repos"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(session *databaseserviceinterface.ServiceOptions, preOutput interface{}) (interface{}, error)
}

type DatabaseServiceTransaction interface {
	RunTransaction(input interface{}) (interface{}, error)
}
