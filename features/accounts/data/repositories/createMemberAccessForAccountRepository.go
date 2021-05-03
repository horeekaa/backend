package accountdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/serviceFailures"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/_exceptionToFailure"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/serviceFailures/enums"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/pkg/errors"
)

type createMemberAccessForAccountRepository struct {
	accountDataSource                            databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource                       databaseaccountdatasourceinterfaces.MemberAccessDataSource
	memberAccessRefDataSource                    databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	createMemberAccessForAccountUsecaseComponent accountdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent
}

func NewCreateMemberAccessForAccountRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (accountdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository, error) {
	return &createMemberAccessForAccountRepository{
		accountDataSource:         accountDataSource,
		memberAccessDataSource:    memberAccessDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
	}, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent,
) (bool, error) {
	createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) preExecute(
	input accountdomainrepositorytypes.CreateMemberAccessForAccountInput,
) (accountdomainrepositorytypes.CreateMemberAccessForAccountInput, error) {
	if &input.Account.ID == nil {
		return accountdomainrepositorytypes.CreateMemberAccessForAccountInput{}, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/createMemberAccessForAccount",
			errors.New(horeekaacorefailureenums.AccountIDNeededToRetrievePersonData),
		)
	}

	if &input.OrganizationType != nil {
		if &input.Organization.ID == nil {
			return accountdomainrepositorytypes.CreateMemberAccessForAccountInput{}, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess,
				"/createMemberAccessForAccount",
				errors.New(horeekaacorefailureenums.OrganizationIDNeededToCreateOrganizationBasedMemberAccess),
			)
		}
	}

	if createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent == nil {
		return input, nil
	}
	return createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent.Validation(input)
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) Execute(
	input accountdomainrepositorytypes.CreateMemberAccessForAccountInput,
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

	memberAccessRef, err := createMbrAccForAccount.memberAccessRefDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"memberAccessRefType":        validatedInput.MemberAccessRefType,
			"organizationType":           validatedInput.OrganizationType,
			"organizationMembershipRole": validatedInput.OrganizationMembershipRole,
			"proposalStatus":             model.EntityProposalStatusApproved,
		},
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
			errors.New(horeekaacorefailureenums.MemberAccessRefNotExist),
		)
	}
	var accessInput model.MemberAccessRefOptionsInput
	jsonTemp, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonTemp, &accessInput)

	memberAccess, err := createMbrAccForAccount.memberAccessDataSource.GetMongoDataSource().Create(
		&model.CreateMemberAccess{
			Account:                    &model.ObjectIDOnly{ID: &account.ID},
			OrganizationMembershipRole: &validatedInput.OrganizationMembershipRole,
			MemberAccessRefType:        validatedInput.MemberAccessRefType,
			Access:                     &accessInput,
			Organization:               &model.ObjectIDOnly{ID: &validatedInput.Organization.ID},
			Status:                     model.MemberAccessStatusActive,
			DefaultAccess:              &model.ObjectIDOnly{ID: &memberAccessRef.ID},
		},
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
