package mongodbcoretypes

import (
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultValuesUpdateType string = "UPDATE"
	DefaultValuesCreateType string = "CREATE"
)

type DefaultValuesOptions struct {
	DefaultValuesType string
}

type OperationOptions struct {
	Session mongodbcorewrapperinterfaces.MongoSessionContext
}

type PaginationOptions struct {
	LastObjectID *primitive.ObjectID
	QueryLimit   *int
}
