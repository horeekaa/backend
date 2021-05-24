package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessForAccountRepository struct {
	accountDataSource                            databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource                       databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	memberAccessRefDataSource                    databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	createMemberAccessForAccountUsecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent
}

func NewCreateMemberAccessForAccountRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository, error) {
	return &createMemberAccessForAccountRepository{
		accountDataSource:         accountDataSource,
		memberAccessDataSource:    memberAccessDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
	}, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent,
) (bool, error) {
	createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) preExecute(
	createMemberAccess *model.CreateMemberAccess,
) (*model.CreateMemberAccess, error) {
	if &createMemberAccess.Account.ID == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/createMemberAccessForAccount",
			nil,
		)
	}

	if createMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased {
		if createMemberAccess.Organization == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess,
				"/createMemberAccessForAccount",
				nil,
			)
		}
	}

	if createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent == nil {
		return createMemberAccess, nil
	}
	return createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent.Validation(createMemberAccess)
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) Execute(
	createMemberAccess *model.CreateMemberAccess,
) (*model.MemberAccess, error) {
	validatedCreateMemberAccess, err := createMbrAccForAccount.preExecute(createMemberAccess)
	if err != nil {
		return nil, err
	}

	account, err := createMbrAccForAccount.accountDataSource.GetMongoDataSource().FindByID(*validatedCreateMemberAccess.Account.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessForAccount",
			err,
		)
	}

	queryMap := map[string]interface{}{
		"memberAccessRefType": validatedCreateMemberAccess.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if validatedCreateMemberAccess.Organization != nil {
		if validatedCreateMemberAccess.Organization.Type != nil {
			queryMap["organizationType"] = *validatedCreateMemberAccess.Organization.Type
		}
	}
	if validatedCreateMemberAccess.OrganizationMembershipRole != nil {
		queryMap["organizationMembershipRole"] = *validatedCreateMemberAccess.OrganizationMembershipRole
	}

	memberAccessRef, err := createMbrAccForAccount.memberAccessRefDataSource.GetMongoDataSource().FindOne(
		queryMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessForAccount",
			err,
		)
	}
	if memberAccessRef == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.MemberAccessRefNotExist,
			"/createMemberAccessForAccount",
			nil,
		)
	}
	var accesscreateMemberAccess model.MemberAccessRefOptionsInput
	jsonTemp, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonTemp, &accesscreateMemberAccess)

	createMemberAccessData := &model.CreateMemberAccess{
		Account:             &model.ObjectIDOnly{ID: &account.ID},
		MemberAccessRefType: validatedCreateMemberAccess.MemberAccessRefType,
		Access:              &accesscreateMemberAccess,
		Status:              model.MemberAccessStatusActive,
		DefaultAccess:       &model.ObjectIDOnly{ID: &memberAccessRef.ID},
	}
	if validatedCreateMemberAccess.Organization != nil {
		createMemberAccessData.Organization = &model.AttachOrganizationInput{
			ID:   validatedCreateMemberAccess.Organization.ID,
			Type: validatedCreateMemberAccess.Organization.Type,
		}
	}
	if validatedCreateMemberAccess.OrganizationMembershipRole != nil {
		createMemberAccessData.OrganizationMembershipRole = validatedCreateMemberAccess.OrganizationMembershipRole
	}

	memberAccess, err := createMbrAccForAccount.memberAccessDataSource.GetMongoDataSource().Create(
		createMemberAccessData,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessForAccount",
			err,
		)
	}

	return memberAccess, nil
}
