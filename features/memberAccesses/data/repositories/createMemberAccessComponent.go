package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessTransactionComponent struct {
	accountDataSource                  databaseaccountdatasourceinterfaces.AccountDataSource
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	memberAccessRefDataSource          databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	memberAccessDataSource             databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	loggingDataSource                  databaseloggingdatasourceinterfaces.LoggingDataSource
	createMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessUsecaseComponent
}

func NewCreateMemberAccessTransactionComponent(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent, error) {
	return &createMemberAccessTransactionComponent{
		accountDataSource:         accountDataSource,
		organizationDataSource:    organizationDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
		memberAccessDataSource:    memberAccessDataSource,
		loggingDataSource:         loggingDataSource,
	}, nil
}

func (createMemberAccessTrx *createMemberAccessTransactionComponent) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessUsecaseComponent,
) (bool, error) {
	createMemberAccessTrx.createMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMemberAccessTrx *createMemberAccessTransactionComponent) PreTransaction(
	input *model.InternalCreateMemberAccess,
) (*model.InternalCreateMemberAccess, error) {
	if input.Account.ID == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/createMemberAccess",
			nil,
		)
	}

	if input.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased {
		if input.Organization == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess,
				"/createMemberAccess",
				nil,
			)
		}
	}

	if createMemberAccessTrx.createMemberAccessUsecaseComponent == nil {
		return input, nil
	}
	return createMemberAccessTrx.createMemberAccessUsecaseComponent.Validation(input)
}

func (createMemberAccessTrx *createMemberAccessTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateMemberAccess,
) (*model.MemberAccess, error) {
	memberAccessToCreate := &model.DatabaseCreateMemberAccess{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, memberAccessToCreate)

	_, err := createMemberAccessTrx.accountDataSource.GetMongoDataSource().FindByID(
		*memberAccessToCreate.Account.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccess",
			err,
		)
	}

	duplicateMemberAccess, err := createMemberAccessTrx.memberAccessDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"account._id":         memberAccessToCreate.Account.ID,
			"status":              model.MemberAccessStatusActive,
			"proposalStatus":      model.EntityProposalStatusApproved,
			"memberAccessRefType": memberAccessToCreate.MemberAccessRefType,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccess",
			err,
		)
	}
	if duplicateMemberAccess != nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.DuplicateObjectExist,
			"/createMemberAccess",
			nil,
		)
	}

	queryMap := map[string]interface{}{
		"memberAccessRefType": memberAccessToCreate.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if memberAccessToCreate.Organization != nil {
		orgToAdd, err := createMemberAccessTrx.organizationDataSource.GetMongoDataSource().FindByID(
			*memberAccessToCreate.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createMemberAccess",
				err,
			)
		}
		queryMap["organizationType"] = orgToAdd.Type

		jsonTemp, _ := json.Marshal(orgToAdd)
		json.Unmarshal(jsonTemp, &memberAccessToCreate.Organization)
	}
	if memberAccessToCreate.OrganizationMembershipRole != nil {
		queryMap["organizationMembershipRole"] = *memberAccessToCreate.OrganizationMembershipRole
	}

	memberAccessRef, err := createMemberAccessTrx.memberAccessRefDataSource.GetMongoDataSource().FindOne(
		queryMap,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccess",
			err,
		)
	}
	if memberAccessRef == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.MemberAccessRefNotExist,
			"/createMemberAccess",
			nil,
		)
	}
	jsonAccessRef, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonAccessRef, &memberAccessToCreate.Access)

	memberAccessToCreate.Status = model.MemberAccessStatusActive

	memberAccessToCreate.DefaultAccessLatestUpdate = &model.ObjectIDOnly{
		ID: &memberAccessRef.ID,
	}

	newDocumentJson, _ := json.Marshal(*memberAccessToCreate)
	generatedObjectID := createMemberAccessTrx.memberAccessDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createMemberAccessTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccess",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: memberAccessToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *memberAccessToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/creatememberAccesse",
			err,
		)
	}

	memberAccessToCreate.ID = generatedObjectID
	memberAccessToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *memberAccessToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		memberAccessToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: memberAccessToCreate.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(memberAccessToCreate)
	json.Unmarshal(jsonTemp, &memberAccessToCreate.ProposedChanges)

	newMemberAccess, err := createMemberAccessTrx.memberAccessDataSource.GetMongoDataSource().Create(
		memberAccessToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccess",
			err,
		)
	}
	return newMemberAccess, nil
}
