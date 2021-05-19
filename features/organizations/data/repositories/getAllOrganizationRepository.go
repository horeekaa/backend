package organizationdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllOrganizationRepository struct {
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewGetAllOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (organizationdomainrepositoryinterfaces.GetAllOrganizationRepository, error) {
	return &getAllOrganizationRepository{
		organizationDataSource,
	}, nil
}

func (getAllMmbAccRefRepo *getAllOrganizationRepository) Execute(
	input organizationdomainrepositorytypes.GetAllOrganizationInput,
) ([]*model.Organization, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

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
