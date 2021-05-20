package organizationdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
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
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

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
