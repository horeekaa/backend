package memberaccessrefdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type getMemberAccessRefRepository struct {
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
}

func NewGetMemberAccessRefRepository(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository, error) {
	return &getMemberAccessRefRepository{
		memberAccessRefDataSource,
	}, nil
}

func (getMmbAccessRefRepo *getMemberAccessRefRepository) Execute(filterFields *model.MemberAccessRefFilterFields) (*model.MemberAccessRef, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := json.Marshal(filterFields)
	json.Unmarshal(data, &filterFieldsMap)

	memberAccessRef, err := getMmbAccessRefRepo.memberAccessRefDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getMemberAccessRef",
			err,
		)
	}

	return memberAccessRef, nil
}
