package loggingdomainrepositories

import (
	"fmt"

	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_exceptionToFailure"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type logEntityApprovalActivity struct {
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	logEntityApprovalActivityUsecaseComponent loggingdomainrepositoryinterfaces.LogEntityApprovalActivityUsecaseComponent
}

func NewLogEntityApprovalActivityRepository(
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository, error) {
	return &logEntityApprovalActivity{
		loggingDataSource: loggingDataSource,
	}, nil
}

func (logEntityAppActivity *logEntityApprovalActivity) SetValidation(
	usecaseComponent loggingdomainrepositoryinterfaces.LogEntityApprovalActivityUsecaseComponent,
) (bool, error) {
	logEntityAppActivity.logEntityApprovalActivityUsecaseComponent = usecaseComponent
	return true, nil
}

func (logEntityAppActivity *logEntityApprovalActivity) preExecute(
	input loggingdomainrepositorytypes.LogEntityApprovalActivityInput,
) (loggingdomainrepositorytypes.LogEntityApprovalActivityInput, error) {
	if logEntityAppActivity.logEntityApprovalActivityUsecaseComponent == nil {
		return input, nil
	}
	return logEntityAppActivity.logEntityApprovalActivityUsecaseComponent.Validation(input)
}

func (logEntityApprovalActivity *logEntityApprovalActivity) Execute(
	input loggingdomainrepositorytypes.LogEntityApprovalActivityInput,
) (*model.Logging, error) {
	validatedInput, err := logEntityApprovalActivity.preExecute(input)
	if err != nil {
		return nil, err
	}

	previousLog, err := logEntityApprovalActivity.loggingDataSource.GetMongoDataSource().FindByID(
		validatedInput.PreviousLog.ID,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/logEntityApprovalActivity",
			err,
		)
	}

	changeLogString := ""
	switch validatedInput.ApprovalStatus {
	case model.EntityProposalStatusRejected:
		changeLogString = fmt.Sprintf(
			"%v rejected the following changes to collection %s:\n%s",
			validatedInput.ApproverInitial,
			previousLog.Collection,
			previousLog.ChangeLog,
		)
		break

	case model.EntityProposalStatusApproved:
		changeLogString = fmt.Sprintf(
			"%v approved the following changes to collection %s:\n%s",
			validatedInput.ApproverInitial,
			previousLog.Collection,
			previousLog.ChangeLog,
		)
		break

	case model.EntityProposalStatusRevisionNeeded:
		changeLogString = fmt.Sprintf(
			"%v asked revision for the following changes to collection %s:\n%s",
			validatedInput.ApproverInitial,
			previousLog.Collection,
			previousLog.ChangeLog,
		)
		break
	}

	createdLog, err := logEntityApprovalActivity.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: previousLog.Collection,
			ChangeLog:  changeLogString,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: &validatedInput.ApprovingAccount.ID,
			},
			Activity:       previousLog.Activity,
			ProposalStatus: validatedInput.ApprovalStatus,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/logEntityApprovalActivity",
			err,
		)
	}

	return createdLog, nil
}
