package accountdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateAccountTransactionComponent struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource  databaseaccountdatasourceinterfaces.PersonDataSource
	pathIdentity      string
}

func NewUpdateAccountTransactionComponent(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (accountdomainrepositoryinterfaces.UpdateAccountTransactionComponent, error) {
	return &updateAccountTransactionComponent{
		accountDataSource: accountDataSource,
		personDataSource:  personDataSource,
		pathIdentity:      "UpdateAccountComponent",
	}, nil
}

func (updateAccountTrx *updateAccountTransactionComponent) PreTransaction(
	updateAccountInput *model.InternalUpdateAccount,
) (*model.InternalUpdateAccount, error) {
	return updateAccountInput, nil
}

func (updateAccountTrx *updateAccountTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateAccount,
) (*model.Account, error) {
	updateAccount := &model.DatabaseUpdateAccount{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateAccount)

	existingAccount, err := updateAccountTrx.accountDataSource.GetMongoDataSource().FindByID(
		updateAccount.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateAccountTrx.pathIdentity,
			err,
		)
	}

	if input.Person != nil {
		existingPerson, err := updateAccountTrx.personDataSource.GetMongoDataSource().FindByID(
			input.Person.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updateAccountTrx.pathIdentity,
				err,
			)
		}
		if existingPerson.ID.Hex() != existingAccount.Person.ID.Hex() {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.InvalidAccountPersonCredential,
				updateAccountTrx.pathIdentity,
				nil,
			)
		}

		_, err = updateAccountTrx.personDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": existingPerson.ID,
			},
			input.Person,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updateAccountTrx.pathIdentity,
				err,
			)
		}
	}

	updatedAccount, err := updateAccountTrx.accountDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": updateAccount.ID,
		},
		updateAccount,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateAccountTrx.pathIdentity,
			err,
		)
	}

	return updatedAccount, nil
}
