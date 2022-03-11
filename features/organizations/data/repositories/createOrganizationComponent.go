package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createOrganizationTransactionComponent struct {
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	loggingDataSource                  databaseloggingdatasourceinterfaces.LoggingDataSource
	createOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent
	generatedObjectID                  *primitive.ObjectID
	pathIdentity                       string
}

func NewCreateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent, error) {
	return &createOrganizationTransactionComponent{
		organizationDataSource: organizationDataSource,
		loggingDataSource:      loggingDataSource,
		pathIdentity:           "CreateOrganizationComponent",
	}, nil
}

func (createOrganizationTrx *createOrganizationTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createOrganizationTrx.organizationDataSource.GetMongoDataSource().GenerateObjectID()
	createOrganizationTrx.generatedObjectID = &generatedObjectID
	return *createOrganizationTrx.generatedObjectID
}

func (createOrganizationTrx *createOrganizationTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createOrganizationTrx.generatedObjectID == nil {
		generatedObjectID := createOrganizationTrx.organizationDataSource.GetMongoDataSource().GenerateObjectID()
		createOrganizationTrx.generatedObjectID = &generatedObjectID
	}
	return *createOrganizationTrx.generatedObjectID
}

func (createOrganizationTrx *createOrganizationTransactionComponent) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent,
) (bool, error) {
	createOrganizationTrx.createOrganizationUsecaseComponent = usecaseComponent
	return true, nil
}

func (createOrganizationTrx *createOrganizationTransactionComponent) PreTransaction(
	input *model.InternalCreateOrganization,
) (*model.InternalCreateOrganization, error) {
	if createOrganizationTrx.createOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return createOrganizationTrx.createOrganizationUsecaseComponent.Validation(input)
}

func (createOrganizationTrx *createOrganizationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateOrganization,
) (*model.Organization, error) {
	organizationToCreate := &model.DatabaseCreateOrganization{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, organizationToCreate)

	newDocumentJson, _ := json.Marshal(*organizationToCreate)
	generatedObjectID := createOrganizationTrx.GetCurrentObjectID()
	loggingOutput, err := createOrganizationTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Organization",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: organizationToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *organizationToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createOrganizationTrx.pathIdentity,
			err,
		)
	}

	organizationToCreate.ID = generatedObjectID
	organizationToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *organizationToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		organizationToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: organizationToCreate.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(organizationToCreate)
	json.Unmarshal(jsonTemp, &organizationToCreate.ProposedChanges)

	newOrganization, err := createOrganizationTrx.organizationDataSource.GetMongoDataSource().Create(
		organizationToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createOrganizationTrx.pathIdentity,
			err,
		)
	}
	createOrganizationTrx.generatedObjectID = nil

	return newOrganization, nil
}
