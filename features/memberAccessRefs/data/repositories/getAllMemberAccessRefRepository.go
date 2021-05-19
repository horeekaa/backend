package memberaccessrefdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllMemberAccessRefRepository struct {
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
}

func NewGetAllMemberAccessRefRepository(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository, error) {
	return &getAllMemberAccessRefRepository{
		memberAccessRefDataSource,
	}, nil
}

func (getAllMmbAccRefRepo *getAllMemberAccessRefRepository) Execute(
	input memberaccessrefdomainrepositorytypes.GetAllMemberAccessRefInput,
) ([]*model.MemberAccessRef, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	memberAccessRefs, err := getAllMmbAccRefRepo.memberAccessRefDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllMemberAccessRef",
			err,
		)
	}

	return memberAccessRefs, nil
}
