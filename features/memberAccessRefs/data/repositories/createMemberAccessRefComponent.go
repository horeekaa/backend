package memberaccessrefdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource             databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	loggingDataSource                     databaseloggingdatasourceinterfaces.LoggingDataSource
	createMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent
}

func NewCreateMemberAccessRefTransactionComponent(
	MemberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent, error) {
	return &createMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: MemberAccessRefDataSource,
		loggingDataSource:         loggingDataSource,
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
	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createMemberAccessRefTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccessRef",
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
