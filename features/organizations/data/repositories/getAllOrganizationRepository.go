package organizationdomainrepositories

import (
	"encoding/json"

	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_exceptionToFailure"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
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
	data, _ := json.Marshal(input.FilterFields)
	json.Unmarshal(data, &filterFieldsMap)

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
