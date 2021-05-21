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
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
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
	input memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput,
) (memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput, error) {
	if &input.Account.ID == nil {
		return memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput{}, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/createMemberAccessForAccount",
			nil,
		)
	}

	if input.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased {
		if input.Organization == nil || input.OrganizationType == "" {
			return memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput{}, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess,
				"/createMemberAccessForAccount",
				nil,
			)
		}
	}

	if createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent == nil {
		return input, nil
	}
	return createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent.Validation(input)
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) Execute(
	input memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput,
) (*model.MemberAccess, error) {
	validatedInput, err := createMbrAccForAccount.preExecute(input)
	if err != nil {
		return nil, err
	}

	account, err := createMbrAccForAccount.accountDataSource.GetMongoDataSource().FindByID(validatedInput.Account.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessForAccount",
			err,
		)
	}

	queryMap := map[string]interface{}{
		"memberAccessRefType": validatedInput.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if validatedInput.OrganizationType != "" {
		queryMap["organizationType"] = validatedInput.OrganizationType
	}
	if validatedInput.OrganizationMembershipRole != "" {
		queryMap["organizationMembershipRole"] = validatedInput.OrganizationMembershipRole
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
	var accessInput model.MemberAccessRefOptionsInput
	jsonTemp, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonTemp, &accessInput)

	createMemberAccessData := &model.CreateMemberAccess{
		Account:             &model.ObjectIDOnly{ID: &account.ID},
		MemberAccessRefType: validatedInput.MemberAccessRefType,
		Access:              &accessInput,
		Status:              model.MemberAccessStatusActive,
		DefaultAccess:       &model.ObjectIDOnly{ID: &memberAccessRef.ID},
	}
	if validatedInput.Organization != nil {
		createMemberAccessData.Organization = &model.ObjectIDOnly{ID: &validatedInput.Organization.ID}
	}
	if validatedInput.OrganizationMembershipRole != "" {
		createMemberAccessData.OrganizationMembershipRole = &validatedInput.OrganizationMembershipRole
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
