package organizationdomainrepositories

import (
	"encoding/json"

	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_exceptionToFailure"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type getOrganizationRepository struct {
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewGetOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (organizationdomainrepositoryinterfaces.GetOrganizationRepository, error) {
	return &getOrganizationRepository{
		organizationDataSource,
	}, nil
}

func (getOrgRepo *getOrganizationRepository) Execute(filterFields *model.OrganizationFilterFields) (*model.Organization, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := json.Marshal(filterFields)
	json.Unmarshal(data, &filterFieldsMap)

	organization, err := getOrgRepo.organizationDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getOrganization",
			err,
		)
	}

	return organization, nil
}
