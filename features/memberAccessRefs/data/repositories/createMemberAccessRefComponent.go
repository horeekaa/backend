package memberaccessrefdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource             databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	loggingDataSource                     databaseloggingdatasourceinterfaces.LoggingDataSource
	structFieldIteratorUtility            coreutilityinterfaces.StructFieldIteratorUtility
	createMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent
}

func NewCreateMemberAccessRefTransactionComponent(
	MemberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent, error) {
	return &createMemberAccessRefTransactionComponent{
		memberAccessRefDataSource:  MemberAccessRefDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
	}, nil
}

func (createMemberAccessRefTrx *createMemberAccessRefTransactionComponent) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent,
) (bool, error) {
	createMemberAccessRefTrx.createMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMemberAccessRefTrx *createMemberAccessRefTransactionComponent) PreTransaction(
	input *model.InternalCreateMemberAccessRef,
) (*model.InternalCreateMemberAccessRef, error) {
	if createMemberAccessRefTrx.createMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return createMemberAccessRefTrx.createMemberAccessRefUsecaseComponent.Validation(input)
}

func (createMemberAccessRefTrx *createMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	fieldChanges := []*model.FieldChangeDataInput{}
	createMemberAccessRefTrx.structFieldIteratorUtility.SetIteratingFunc(
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
	createMemberAccessRefTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createMemberAccessRefTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	generatedObjectID := createMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createMemberAccessRefTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccessRef",
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
			"/createMemberAccessRef",
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

	newMemberAccessRef, err := createMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().Create(
		input,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessRef",
			err,
		)
	}
	return newMemberAccessRef, nil
}
