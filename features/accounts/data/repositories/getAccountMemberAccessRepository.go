package accountdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAccountMemberAccessRepository struct {
	accountDataSource                      databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource                 databaseaccountdatasourceinterfaces.MemberAccessDataSource
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
	getAccountMemberAccessUsecaseComponent accountdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent
}

func NewGetAccountMemberAccessRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository, error) {
	return &getAccountMemberAccessRepository{
		accountDataSource:      accountDataSource,
		memberAccessDataSource: memberAccessDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent,
) (bool, error) {
	getAccountMemberAccess.getAccountMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) preExecute(
	input accountdomainrepositorytypes.GetAccountMemberAccessInput,
) (accountdomainrepositorytypes.GetAccountMemberAccessInput, error) {
	if &input.Account.ID == nil {
		return accountdomainrepositorytypes.GetAccountMemberAccessInput{}, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/getAccountMemberAccess",
			nil,
		)
	}
	if getAccountMemberAccess.getAccountMemberAccessUsecaseComponent == nil {
		return input, nil
	}
	return getAccountMemberAccess.getAccountMemberAccessUsecaseComponent.Validation(input)
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) Execute(input accountdomainrepositorytypes.GetAccountMemberAccessInput) (*model.MemberAccess, error) {
	preExecuteOutput, err := getAccountMemberAccess.preExecute(input)
	if err != nil {
		return nil, err
	}

	account, err := getAccountMemberAccess.accountDataSource.GetMongoDataSource().FindByID(
		preExecuteOutput.Account.ID,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAccountMemberAccess",
			err,
		)
	}
	var accessMap map[string]interface{}
	jsonTemp, _ := json.Marshal(preExecuteOutput.MemberAccessRefOptions)
	json.Unmarshal(jsonTemp, &accessMap)

	getMemberAccessQuery := make(map[string]interface{})
	getAccountMemberAccess.mapProcessorUtility.FlattenMap(
		"",
		map[string]interface{}{
			"account":             map[string]interface{}{"_id": account.ID},
			"memberAccessRefType": preExecuteOutput.MemberAccessRefType,
			"access":              accessMap,
			"status":              model.MemberAccessStatusActive,
		},
		&getMemberAccessQuery,
	)

	memberAccess, err := getAccountMemberAccess.memberAccessDataSource.GetMongoDataSource().FindOne(
		getMemberAccessQuery,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAccountMemberAccess",
			err,
		)
	}
	if memberAccess == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.FeatureNotAccessibleByAccount,
			"/getAccountMemberAccess",
			nil,
		)
	}
	return memberAccess, nil
}
