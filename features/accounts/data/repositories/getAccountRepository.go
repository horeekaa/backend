package accountdomainrepositories

import (
	"encoding/json"

	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_exceptionToFailure"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type getAccountRepository struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
}

func NewGetAccountRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.GetAccountRepository, error) {
	return &getAccountRepository{
		accountDataSource,
	}, nil
}

func (getMmbAccessRefRepo *getAccountRepository) Execute(filterFields *model.AccountFilterFields) (*model.Account, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := json.Marshal(filterFields)
	json.Unmarshal(data, &filterFieldsMap)

	account, err := getMmbAccessRefRepo.accountDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAccount",
			err,
		)
	}

	return account, nil
}
