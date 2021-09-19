package organizationdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllOrganizationRepository struct {
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
	mongoQueryBuilder      mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (organizationdomainrepositoryinterfaces.GetAllOrganizationRepository, error) {
	return &getAllOrganizationRepository{
		organizationDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllMmbAccRefRepo *getAllOrganizationRepository) Execute(
	input organizationdomainrepositorytypes.GetAllOrganizationInput,
) ([]*model.Organization, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMmbAccRefRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	organizations, err := getAllMmbAccRefRepo.organizationDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllOrganization",
			err,
		)
	}

	return organizations, nil
}
