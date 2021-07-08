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
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessForAccountRepository struct {
	accountDataSource                            databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource                       databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	memberAccessRefDataSource                    databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	organizationDataSource                       databaseorganizationdatasourceinterfaces.OrganizationDataSource
	createMemberAccessForAccountUsecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent
}

func NewCreateMemberAccessForAccountRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository, error) {
	return &createMemberAccessForAccountRepository{
		accountDataSource:         accountDataSource,
		memberAccessDataSource:    memberAccessDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
		organizationDataSource:    organizationDataSource,
	}, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountUsecaseComponent,
) (bool, error) {
	createMbrAccForAccount.createMemberAccessForAccountUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMbrAccForAccount *createMemberAccessForAccountRepository) preExecute(
	createMemberAccess *model.InternalCreateMemberAccess,
) (*model.InternalCreateMemberAccess, error) {
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
	createMemberAccess *model.InternalCreateMemberAccess,
) (*model.MemberAccess, error) {
	validatedCreateMemberAccess, err := createMbrAccForAccount.preExecute(createMemberAccess)
	if err != nil {
		return nil, err
	}

	_, err = createMbrAccForAccount.accountDataSource.GetMongoDataSource().FindByID(*validatedCreateMemberAccess.Account.ID, &mongodbcoretypes.OperationOptions{})
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
		orgToAdd, err := createMbrAccForAccount.organizationDataSource.GetMongoDataSource().FindByID(
			validatedCreateMemberAccess.Organization.ID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createMemberAccess",
				err,
			)
		}

		jsonTemp, _ := json.Marshal(orgToAdd)
		json.Unmarshal(jsonTemp, &validatedCreateMemberAccess.Organization)
		json.Unmarshal(jsonTemp, &validatedCreateMemberAccess.OrganizationLatestUpdate)
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
	jsonTemp, _ := json.Marshal(memberAccessRef.Access)
	json.Unmarshal(jsonTemp, &validatedCreateMemberAccess.Access)

	validatedCreateMemberAccess.Status = model.MemberAccessStatusActive

	jsonTemp, _ = json.Marshal(memberAccessRef)
	json.Unmarshal(jsonTemp, &validatedCreateMemberAccess.DefaultAccess)
	validatedCreateMemberAccess.DefaultAccessLatestUpdate = &model.ObjectIDOnly{
		ID: &memberAccessRef.ID,
	}

	jsonTemp, _ = json.Marshal(validatedCreateMemberAccess)
	json.Unmarshal(jsonTemp, &validatedCreateMemberAccess.ProposedChanges)

	newMemberAccess, err := createMbrAccForAccount.memberAccessDataSource.GetMongoDataSource().Create(
		validatedCreateMemberAccess,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccess",
			err,
		)
	}
	return newMemberAccess, nil
}
