package accountdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAccountRepository struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	pathIdentity      string
}

func NewGetAccountRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.GetAccountRepository, error) {
	return &getAccountRepository{
		accountDataSource,
		"GetAccountFromAuthRepo",
	}, nil
}

func (getMmbAccessRefRepo *getAccountRepository) Execute(filterFields *model.AccountFilterFields) (*model.Account, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	account, err := getMmbAccessRefRepo.accountDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getMmbAccessRefRepo.pathIdentity,
			err,
		)
	}

	return account, nil
}
