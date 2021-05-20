package loggingdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getLoggingRepository struct {
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource
}

func NewGetLoggingRepository(
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (loggingdomainrepositoryinterfaces.GetLoggingRepository, error) {
	return &getLoggingRepository{
		loggingDataSource,
	}, nil
}

func (getLoggingRefRepo *getLoggingRepository) Execute(filterFields *model.LoggingFilterFields) (*model.Logging, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	logging, err := getLoggingRefRepo.loggingDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getLogging",
			err,
		)
	}

	return logging, nil
}
