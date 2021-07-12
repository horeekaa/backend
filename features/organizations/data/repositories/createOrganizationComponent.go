package organizationdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createOrganizationTransactionComponent struct {
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	loggingDataSource                  databaseloggingdatasourceinterfaces.LoggingDataSource
	structFieldIteratorUtility         coreutilityinterfaces.StructFieldIteratorUtility
	createOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent
}

func NewCreateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent, error) {
	return &createOrganizationTransactionComponent{
		organizationDataSource:     organizationDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
	}, nil
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
	fieldChanges := []*model.FieldChangeDataInput{}
	createOrganizationTrx.structFieldIteratorUtility.SetIteratingFunc(
		func(tag interface{}, field interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field).Kind().String(),
				NewValue: fmt.Sprint(field),
			})
			*tagString = ""
		},
	)
	createOrganizationTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createOrganizationTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	generatedObjectID := createOrganizationTrx.organizationDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createOrganizationTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Organization",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			FieldChanges: fieldChanges,
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

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, &input.ProposedChanges)

	newOrganization, err := createOrganizationTrx.organizationDataSource.GetMongoDataSource().Create(
		input,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createOrganization",
			err,
		)
	}
	return newOrganization, nil
}
