package memberaccessrefdomainrepositories

import (
	"encoding/json"
	"time"

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
	pathIdentity                          string
}

func NewCreateMemberAccessRefTransactionComponent(
	MemberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent, error) {
	return &createMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: MemberAccessRefDataSource,
		loggingDataSource:         loggingDataSource,
		pathIdentity:              "CreateMemberAccessRefComponent",
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
	memberAccessRefToCreate := &model.DatabaseCreateMemberAccessRef{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, memberAccessRefToCreate)

	newDocumentJson, _ := json.Marshal(*memberAccessRefToCreate)
	generatedObjectID := createMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createMemberAccessRefTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccessRef",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: memberAccessRefToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *memberAccessRefToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createMemberAccessRefTrx.pathIdentity,
			err,
		)
	}

	memberAccessRefToCreate.ID = generatedObjectID
	memberAccessRefToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *memberAccessRefToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		memberAccessRefToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: memberAccessRefToCreate.SubmittingAccount.ID}
	}
	var currentTime = time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed

	if memberAccessRefToCreate.ProposalStatus == nil {
		memberAccessRefToCreate.ProposalStatus = &defaultProposalStatus
	}
	memberAccessRefToCreate.CreatedAt = &currentTime
	memberAccessRefToCreate.UpdatedAt = &currentTime

	jsonTemp, _ = json.Marshal(memberAccessRefToCreate)
	json.Unmarshal(jsonTemp, &memberAccessRefToCreate.ProposedChanges)

	newMemberAccessRef, err := createMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().Create(
		memberAccessRefToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createMemberAccessRefTrx.pathIdentity,
			err,
		)
	}
	return newMemberAccessRef, nil
}
