package mongodbcoreclientinterfaces

import (
	"time"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
)

type MongoClient interface {
	Connect() (bool, error)
	GetDatabaseName() (string, error)
	GetDatabaseTimeout() (time.Duration, error)
	GetCollectionRef(collectionName string) (mongodbcorewrapperinterfaces.MongoCollectionRef, error)
	CreateNewSession() (mongodbcorewrapperinterfaces.MongoSession, error)
}
