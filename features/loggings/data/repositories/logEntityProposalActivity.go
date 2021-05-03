package loggingdomainrepositories

import (
	"errors"
	"fmt"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/serviceFailures"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/_exceptionToFailure"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/serviceFailures/enums"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type logEntityProposalActivityRepository struct {
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	structComparisonUtility                   coreutilityinterfaces.StructComparisonUtility
	structFieldIteratorUtility                coreutilityinterfaces.StructFieldIteratorUtility
	logEntityProposalActivityUsecaseComponent loggingdomainrepositoryinterfaces.LogEntityProposalActivityUsecaseComponent
}

func NewLogEntityProposalActivityRepository(
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository, error) {
	return &logEntityProposalActivityRepository{
		loggingDataSource:          loggingDataSource,
		structComparisonUtility:    structComparisonUtility,
		structFieldIteratorUtility: structFieldIteratorUtility,
	}, nil
}

func (logEttPrpsalActivity *logEntityProposalActivityRepository) SetValidation(
	usecaseComponent loggingdomainrepositoryinterfaces.LogEntityProposalActivityUsecaseComponent,
) (bool, error) {
	logEttPrpsalActivity.logEntityProposalActivityUsecaseComponent = usecaseComponent
	return true, nil
}

func (logEttPrpsalActivity *logEntityProposalActivityRepository) preExecute(
	input loggingdomainrepositorytypes.LogEntityProposalActivityInput,
) (loggingdomainrepositorytypes.LogEntityProposalActivityInput, error) {
	if logEttPrpsalActivity.logEntityProposalActivityUsecaseComponent == nil {
		return input, nil
	}
	return logEttPrpsalActivity.logEntityProposalActivityUsecaseComponent.Validation(input)
}

func (logEttPrpsalActivity *logEntityProposalActivityRepository) Execute(
	input loggingdomainrepositorytypes.LogEntityProposalActivityInput,
) (*model.Logging, error) {
	validatedInput, err := logEttPrpsalActivity.preExecute(input)
	if err != nil {
		return nil, err
	}
	var changeLogString interface{} = ""
	creatorInitial := fmt.Sprintf("%s ", validatedInput.CreatorInitial)
	if validatedInput.ProposalStatus == model.EntityProposalStatusApproved {
		changeLogString = fmt.Sprintf(
			"%s approved the ",
			creatorInitial,
		)
		creatorInitial = "self-"
	}

	switch validatedInput.Activity {
	case model.LoggedActivityUpdate:
		if validatedInput.ExistingObject == nil || validatedInput.ExistingObjectID == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity,
				"/logEntityProposalActivity",
				errors.New(horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity),
			)
		}
		existingIDValue := *validatedInput.ExistingObjectID
		changeLogString = fmt.Sprintf(
			"%s %sproposed changes to %s with id **%s to the followings:\n",
			changeLogString,
			creatorInitial,
			validatedInput.CollectionName,
			string(existingIDValue[len(existingIDValue)-5:len(existingIDValue)-1]),
		)

		logEttPrpsalActivity.structComparisonUtility.SetComparisonFunc(
			func(tag interface{}, field1 interface{}, field2 interface{}, logString *interface{}) {
				if field1 == field2 {
					return
				}
				*logString = fmt.Sprintf(
					"%v%v from %v to %v\n",
					*logString,
					tag,
					field2,
					field1,
				)
			},
		)
		logEttPrpsalActivity.structComparisonUtility.SetPreDeepComparisonFunc(
			func(tag interface{}, logString *interface{}) {
				*logString = fmt.Sprintf(
					"%v%v.",
					*logString,
					tag,
				)
			},
		)
		logEttPrpsalActivity.structComparisonUtility.CompareStructs(
			*validatedInput.NewObject,
			*validatedInput.ExistingObject,
			&changeLogString,
		)
		break

	case model.LoggedActivityCreate:
		changeLogString = fmt.Sprintf(
			"%s %sproposed creation to %s with the followings:\n",
			changeLogString,
			creatorInitial,
			validatedInput.CollectionName,
		)
		logEttPrpsalActivity.structFieldIteratorUtility.SetIteratingFunc(
			func(tag interface{}, field interface{}, logString *interface{}) {
				*logString = fmt.Sprintf(
					"%v%v with value %v\n",
					*logString,
					tag,
					field,
				)
			},
		)
		logEttPrpsalActivity.structFieldIteratorUtility.SetPreDeepIterateFunc(
			func(tag interface{}, logString *interface{}) {
				*logString = fmt.Sprintf(
					"%v%v.",
					*logString,
					tag,
				)
			},
		)
		logEttPrpsalActivity.structFieldIteratorUtility.IterateStruct(
			*validatedInput.NewObject,
			&changeLogString,
		)
		break

	case model.LoggedActivityDelete:
		if validatedInput.ExistingObject == nil || validatedInput.ExistingObjectID == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity,
				"/logEntityProposalActivity",
				errors.New(horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity),
			)
		}
		existingIDValue := *validatedInput.ExistingObjectID
		changeLogString = fmt.Sprintf(
			"%s %sproposed deletion to %s with id **%s\n",
			changeLogString,
			creatorInitial,
			validatedInput.CollectionName,
			string(existingIDValue[len(existingIDValue)-5:len(existingIDValue)-1]),
		)
		break
	}

	loggingOutput, err := logEttPrpsalActivity.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: validatedInput.CollectionName,
			ChangeLog:  changeLogString.(string),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: &validatedInput.CreatedByAccount.ID,
			},
			Activity:       validatedInput.Activity,
			ProposalStatus: validatedInput.ProposalStatus,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/logEntityProposalActivity",
			err,
		)
	}

	return loggingOutput, nil
}
