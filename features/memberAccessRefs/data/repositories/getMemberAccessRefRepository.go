package memberaccessrefdomainrepositories

import (
	"encoding/json"

	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_exceptionToFailure"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberaccessrefs/data/dataSources/databases/interfaces/sources"
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

func (getMmbAccessRefRepo *getMemberAccessRefRepository) Execute(filterFields *model.UpdateMemberAccessRef) (*model.MemberAccessRef, error) {
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
