package memberaccessdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
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
	structFieldIteratorUtility         coreutilityinterfaces.StructFieldIteratorUtility
	createMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessUsecaseComponent
}

func NewCreateMemberAccessTransactionComponent(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent, error) {
	return &createMemberAccessTransactionComponent{
		accountDataSource:          accountDataSource,
		organizationDataSource:     organizationDataSource,
		memberAccessRefDataSource:  memberAccessRefDataSource,
		memberAccessDataSource:     memberAccessDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
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
	_, err := createMemberAccessTrx.accountDataSource.GetMongoDataSource().FindByID(
		*input.Account.ID,
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
			"account._id":         input.Account.ID,
			"status":              model.MemberAccessStatusActive,
			"proposalStatus":      model.EntityProposalStatusApproved,
			"memberAccessRefType": input.MemberAccessRefType,
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

	fieldChanges := []*model.FieldChangeDataInput{}
	createMemberAccessTrx.structFieldIteratorUtility.SetIteratingFunc(
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
	createMemberAccessTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createMemberAccessTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	queryMap := map[string]interface{}{
		"memberAccessRefType": input.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if input.Organization != nil {
		orgToAdd, err := createMemberAccessTrx.organizationDataSource.GetMongoDataSource().FindByID(
			input.Organization.ID,
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
		json.Unmarshal(jsonTemp, &input.Organization)
		json.Unmarshal(jsonTemp, &input.OrganizationLatestUpdate)
	}
	if input.OrganizationMembershipRole != nil {
		queryMap["organizationMembershipRole"] = *input.OrganizationMembershipRole
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
	jsonTemp, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonTemp, &input.Access)

	input.Status = model.MemberAccessStatusActive

	jsonTemp, _ = json.Marshal(memberAccessRef)
	json.Unmarshal(jsonTemp, &input.DefaultAccess)
	input.DefaultAccessLatestUpdate = &model.ObjectIDOnly{
		ID: &memberAccessRef.ID,
	}

	generatedObjectID := createMemberAccessTrx.memberAccessDataSource.GetMongoDataSource().GenerateObjectID()
	loggingOutput, err := createMemberAccessTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccess",
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
			"/creatememberAccesse",
			err,
		)
	}

	input.ID = generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(input)
	json.Unmarshal(jsonTemp, &input.ProposedChanges)

	newMemberAccess, err := createMemberAccessTrx.memberAccessDataSource.GetMongoDataSource().Create(
		input,
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
