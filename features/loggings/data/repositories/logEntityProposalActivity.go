package loggingdomainrepositories

import (
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
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
	fieldChanges := []*model.FieldChangeDataInput{}

	switch validatedInput.Activity {
	case model.LoggedActivityUpdate:
		if validatedInput.ExistingObject == nil || validatedInput.ExistingObjectID == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity,
				"/logEntityProposalActivity",
				nil,
			)
		}

		logEttPrpsalActivity.structComparisonUtility.SetComparisonFunc(
			func(tag interface{}, field1 interface{}, field2 interface{}, tagString *interface{}) {
				if field1 == field2 {
					return
				}
				*tagString = fmt.Sprintf(
					"%v%v",
					*tagString,
					tag,
				)

				fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
					Name:     fmt.Sprint(tagString),
					Type:     reflect.TypeOf(field1).Kind().String(),
					OldValue: fmt.Sprint(field2),
					NewValue: fmt.Sprint(field1),
				})
				*tagString = ""
			},
		)
		logEttPrpsalActivity.structComparisonUtility.SetPreDeepComparisonFunc(
			func(tag interface{}, tagString *interface{}) {
				*tagString = fmt.Sprintf(
					"%v%v.",
					*tagString,
					tag,
				)
			},
		)
		var tagString interface{} = ""
		logEttPrpsalActivity.structComparisonUtility.CompareStructs(
			*validatedInput.NewObject,
			*validatedInput.ExistingObject,
			&tagString,
		)
		break

	case model.LoggedActivityCreate:
		logEttPrpsalActivity.structFieldIteratorUtility.SetIteratingFunc(
			func(tag interface{}, field interface{}, tagString *interface{}) {
				*tagString = fmt.Sprintf(
					"%v%v",
					*tagString,
					tag,
				)

				fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
					Name:     fmt.Sprint(tagString),
					Type:     reflect.TypeOf(field).Kind().String(),
					NewValue: fmt.Sprint(field),
				})
				*tagString = ""
			},
		)
		logEttPrpsalActivity.structFieldIteratorUtility.SetPreDeepIterateFunc(
			func(tag interface{}, tagString *interface{}) {
				*tagString = fmt.Sprintf(
					"%v%v.",
					*tagString,
					tag,
				)
			},
		)
		var tagString interface{} = ""
		logEttPrpsalActivity.structFieldIteratorUtility.IterateStruct(
			*validatedInput.NewObject,
			&tagString,
		)
		break

	case model.LoggedActivityDelete:
		if validatedInput.ExistingObject == nil || validatedInput.ExistingObjectID == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.ExistingObjectAndItsIDMustNotBeNilForUpdateActivity,
				"/logEntityProposalActivity",
				nil,
			)
		}
		break
	}

	loggingOutput, err := logEttPrpsalActivity.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: validatedInput.CollectionName,
			Document:   validatedInput.ExistingObjectID,
			FieldChanges:  fieldChanges,
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
