package databaseserviceinterface

import (
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type ServiceOptions struct {
	OperationOptions *mongooperations.OperationOptions
}
