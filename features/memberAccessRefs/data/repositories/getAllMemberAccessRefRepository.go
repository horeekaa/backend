package memberaccessrefdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllMemberAccessRefRepository struct {
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	mongoQueryBuilder         mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity              string
}

func NewGetAllMemberAccessRefRepository(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository, error) {
	return &getAllMemberAccessRefRepository{
		memberAccessRefDataSource,
		mongoQueryBuilder,
		"GetAllMemberAccessRefRepository",
	}, nil
}

func (getAllMmbAccRefRepo *getAllMemberAccessRefRepository) Execute(
	input memberaccessrefdomainrepositorytypes.GetAllMemberAccessRefInput,
) ([]*model.MemberAccessRef, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMmbAccRefRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	memberAccessRefs, err := getAllMmbAccRefRepo.memberAccessRefDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllMmbAccRefRepo.pathIdentity,
			err,
		)
	}

	return memberAccessRefs, nil
}
