package memberaccessrefpresentationusecases

import (
	"errors"
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/_errors/usecaseErrors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/_errors/usecaseErrors/_enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefUsecase struct {
	manageAccountAuthenticationRepo     accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository
	getAccountMemberAccessRepo          accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo        accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	createMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository
	logEntityProposalActivityRepo       loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	createMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateMemberAccessRefUsecase(
	manageAccountAuthenticationRepo accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	createMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
) (memberaccessrefpresentationusecaseinterfaces.CreateMemberAccessRefUsecase, error) {
	return &createMemberAccessRefUsecase{
		manageAccountAuthenticationRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		createMemberAccessRefRepo,
		logEntityProposalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			MemberAccessRefAccesses: &model.MemberAccessRefAccessesInput{
				MemberAccessRefCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMmbAccessRefUcase *createMemberAccessRefUsecase) validation(input memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput) (memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput, error) {
	if &input.AuthHeader == nil {
		return memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"createMemberAccessRefUsecase/",
				errors.New("createMemberAccessRefUsecase/"),
			)
	}
	return input, nil
}

func (createMmbAccessRefUcase *createMemberAccessRefUsecase) Execute(input memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput) (*model.MemberAccessRef, error) {
	validatedInput, err := createMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMmbAccessRefUcase.manageAccountAuthenticationRepo.RunTransaction(
		accountdomainrepositorytypes.ManageAccountAuthenticationInput{
			AuthHeader: validatedInput.AuthHeader,
			Context:    validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"createMemberAccessRefUsecase/",
			err,
		)
	}

	personChannel := make(chan *model.Person)
	errChannel := make(chan error)
	go func() {
		person, err := createMmbAccessRefUcase.getPersonDataFromAccountRepo.Execute(account)
		if err != nil {
			errChannel <- err
		}
		personChannel <- person
	}()

	accMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		accountdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *createMmbAccessRefUcase.createMemberAccessRefAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"createMemberAccessRefUsecase/",
			err,
		)
	}
	if *accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval {
		validatedInput.CreateMemberAccessRef.ProposalStatus =
			func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
	}

	createdMemberAccessRef, err := createMmbAccessRefUcase.createMemberAccessRefRepo.Execute(
		validatedInput.CreateMemberAccessRef,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"createMemberAccessRefUsecase/",
			err,
		)
	}

	accountInitials := ""
	select {
	case person := <-personChannel:
		accountInitials = fmt.Sprintf("XXXX%s", account.ID.Hex()[len(account.ID.Hex())-6:])
		if person != nil {
			accountInitials = person.FirstName
		}

		break
	case err := <-errChannel:
		return nil, err
	}

	var newObject interface{} = *createdMemberAccessRef
	_, err = createMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "MemberAccessRef",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityCreate,
			ProposalStatus:   *validatedInput.CreateMemberAccessRef.ProposalStatus,
			NewObject:        &newObject,
			CreatorInitial:   accountInitials,
		},
	)

	return createdMemberAccessRef, nil
}
