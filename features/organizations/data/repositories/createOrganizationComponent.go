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
}

func NewCreateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent, error) {
	return &createOrganizationTransactionComponent{
		organizationDataSource: organizationDataSource,
		loggingDataSource:      loggingDataSource,
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
	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createOrganizationTrx.GetCurrentObjectID()
	loggingOutput, err := createOrganizationTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Organization",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: input.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *input.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createOrganization",
			err,
		)
	}

	input.ID = generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	organizationToCreate := &model.DatabaseCreateOrganization{}

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, organizationToCreate)
	json.Unmarshal(jsonTemp, &organizationToCreate.ProposedChanges)

	newOrganization, err := createOrganizationTrx.organizationDataSource.GetMongoDataSource().Create(
		organizationToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createOrganization",
			err,
		)
	}
	createOrganizationTrx.generatedObjectID = nil

	return newOrganization, nil
}
