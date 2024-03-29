package memberaccessdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
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
	mongoQueryBuilder      mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity           string
}

func NewGetAllMemberAccessRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository, error) {
	return &getAllMemberAccessRepository{
		memberAccessDataSource,
		mongoQueryBuilder,
		"GetAllMemberAccessRepository",
	}, nil
}

func (getAllMemberAccessRepo *getAllMemberAccessRepository) Execute(
	input memberaccessdomainrepositorytypes.GetAllMemberAccessInput,
) ([]*model.MemberAccess, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMemberAccessRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	var mongoPagination mongodbcoretypes.PaginationOptions
	data, _ := bson.Marshal(input.PaginationOpt)
	bson.Unmarshal(data, &mongoPagination)

	memberAccesses, err := getAllMemberAccessRepo.memberAccessDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllMemberAccessRepo.pathIdentity,
			err,
		)
	}

	return memberAccesses, nil
}
