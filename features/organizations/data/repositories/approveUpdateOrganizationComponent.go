package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateOrganizationTransactionComponent struct {
	organizationDataSource                    databaseorganizationdatasourceinterfaces.OrganizationDataSource
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                       coreutilityinterfaces.MapProcessorUtility
	approveUpdateOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationUsecaseComponent
}

func NewApproveUpdateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent, error) {
	return &approveUpdateOrganizationTransactionComponent{
		organizationDataSource: organizationDataSource,
		loggingDataSource:      loggingDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (approveProdTrx *approveUpdateOrganizationTransactionComponent) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationUsecaseComponent,
) (bool, error) {
	approveProdTrx.approveUpdateOrganizationUsecaseComponent = usecaseComponent
	return true, nil
}

func (approveProdTrx *approveUpdateOrganizationTransactionComponent) PreTransaction(
	input *model.InternalUpdateOrganization,
) (*model.InternalUpdateOrganization, error) {
	if approveProdTrx.approveUpdateOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return approveProdTrx.approveUpdateOrganizationUsecaseComponent.Validation(input)
}

func (approveProdTrx *approveUpdateOrganizationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateOrganization *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	existingOrganization, err := approveProdTrx.organizationDataSource.GetMongoDataSource().FindByID(
		updateOrganization.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}

	previousLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingOrganization.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateOrganization.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateOrganization.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"FieldChanges": previousLog.FieldChanges,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateOrganization.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateOrganization := &model.InternalUpdateOrganization{
		ID: updateOrganization.ID,
	}
	jsonExisting, _ := json.Marshal(existingOrganization.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateOrganization.ProposedChanges)

	var updateOrganizationMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateOrganization)
	json.Unmarshal(jsonUpdate, &updateOrganizationMap)

	approveProdTrx.mapProcessorUtility.RemoveNil(updateOrganizationMap)

	jsonUpdate, _ = json.Marshal(updateOrganizationMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateOrganization.ProposedChanges)

	if updateOrganization.ProposalStatus != nil {
		if *updateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateOrganization.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateOrganization)
		}
	}

	updatedOrganization, err := approveProdTrx.organizationDataSource.GetMongoDataSource().Update(
		fieldsToUpdateOrganization.ID,
		fieldsToUpdateOrganization,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}

	return updatedOrganization, nil
}
