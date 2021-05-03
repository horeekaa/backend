package memberaccessrefdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/_exceptionToFailure"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefRepository struct {
	memberAccessRefDataSource             databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	createMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent
}

func NewCreateMemberAccessRefRepository(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository, error) {
	return &createMemberAccessRefRepository{
		memberAccessRefDataSource: memberAccessRefDataSource,
	}, nil
}

func (createMmbAccRefRepo *createMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent,
) (bool, error) {
	createMmbAccRefRepo.createMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (createMmbAccRefRepo *createMemberAccessRefRepository) preExecute(
	input *model.CreateMemberAccessRef,
) (*model.CreateMemberAccessRef, error) {
	if createMmbAccRefRepo.createMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return createMmbAccRefRepo.createMemberAccessRefUsecaseComponent.Validation(input)
}

func (createMmbAccRefRepo *createMemberAccessRefRepository) Execute(
	input *model.CreateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	validatedInput, err := createMmbAccRefRepo.preExecute(input)
	if err != nil {
		return nil, err
	}

	newMemberMmbRefRepo, err := createMmbAccRefRepo.memberAccessRefDataSource.GetMongoDataSource().Create(
		validatedInput,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMemberAccessRef",
			err,
		)
	}
	return newMemberMmbRefRepo, nil
}
