package memberaccessdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllMemberAccessRepository struct {
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource
}

func NewGetAllMemberAccessRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
) (memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository, error) {
	return &getAllMemberAccessRepository{
		memberAccessDataSource,
	}, nil
}

func (getAllmmbAccRepo *getAllMemberAccessRepository) Execute(
	input memberaccessdomainrepositorytypes.GetAllMemberAccessInput,
) ([]*model.MemberAccess, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	var mongoPagination mongodbcoretypes.PaginationOptions
	data, _ = bson.Marshal(input.PaginationOpt)
	bson.Unmarshal(data, &mongoPagination)

	memberAccesses, err := getAllmmbAccRepo.memberAccessDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllmemberAccess",
			err,
		)
	}

	return memberAccesses, nil
}
