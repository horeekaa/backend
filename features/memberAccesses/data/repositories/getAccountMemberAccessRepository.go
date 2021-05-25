package memberaccessdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAccountMemberAccessRepository struct {
	memberAccessDataSource                 databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
	getAccountMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent
}

func NewGetAccountMemberAccessRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository, error) {
	return &getAccountMemberAccessRepository{
		memberAccessDataSource: memberAccessDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent,
) (bool, error) {
	getAccountMemberAccess.getAccountMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) preExecute(
	getAccMmbAccInput memberaccessdomainrepositorytypes.GetAccountMemberAccessInput,
) (memberaccessdomainrepositorytypes.GetAccountMemberAccessInput, error) {
	if !getAccMmbAccInput.QueryMode {
		if &getAccMmbAccInput.MemberAccessFilterFields.Account == nil {
			return memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{}, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
				"/getAccountMemberAccess",
				nil,
			)
		}
	}
	if getAccountMemberAccess.getAccountMemberAccessUsecaseComponent == nil {
		return getAccMmbAccInput, nil
	}
	return getAccountMemberAccess.getAccountMemberAccessUsecaseComponent.Validation(getAccMmbAccInput)
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) Execute(getMmbAccInput memberaccessdomainrepositorytypes.GetAccountMemberAccessInput) (*model.MemberAccess, error) {
	validatedInput, err := getAccountMemberAccess.preExecute(getMmbAccInput)
	if err != nil {
		return nil, err
	}
	if validatedInput.MemberAccessFilterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(validatedInput.MemberAccessFilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	getMemberAccessQuery := make(map[string]interface{})
	getAccountMemberAccess.mapProcessorUtility.FlattenMap(
		"",
		filterFieldsMap,
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
	if memberAccess == nil && !validatedInput.QueryMode {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.FeatureNotAccessibleByAccount,
			"/getAccountMemberAccess",
			nil,
		)
	}
	return memberAccess, nil
}
