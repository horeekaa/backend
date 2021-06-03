package memberaccesspresentationusecases

import (
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepo     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	createMemberAccessRepo           memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository
	logEntityProposalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	createMemberAccessAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateMemberAccessUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
) (memberaccesspresentationusecaseinterfaces.CreateMemberAccessUsecase, error) {
	return &createMemberAccessUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		createMemberAccessRepo,
		logEntityProposalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMmbAccessRefUcase *createMemberAccessUsecase) validation(input memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput) (memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput, error) {
	if &input.Context == nil {
		return memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createMemberAccessUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateMemberAccess.SubmittingAccount = nil
	input.CreateMemberAccess.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createMmbAccessRefUcase *createMemberAccessUsecase) Execute(input memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput) (*model.MemberAccess, error) {
	validatedInput, err := createMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMmbAccessRefUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createMemberAccessUsecase",
			nil,
		)
	}

	duplicateMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account: &model.ObjectIDOnly{ID: validatedInput.CreateMemberAccess.Account.ID},
				Status: func(m model.MemberAccessStatus) *model.MemberAccessStatus {
					return &m
				}(model.MemberAccessStatusActive),
				ProposalStatus: func(m model.EntityProposalStatus) *model.EntityProposalStatus {
					return &m
				}(model.EntityProposalStatusApproved),
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}
	if duplicateMemberAccess != nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.DuplicateAccessExist,
			403,
			"/createMemberAccessUsecase",
			nil,
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

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Access:              createMmbAccessRefUcase.createMemberAccessAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}
	if accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval != nil {
		if *accMemberAccess.Access.ManageMemberAccesses.MemberAccessApproval {
			validatedInput.CreateMemberAccess.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
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
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}

	var newObject interface{} = *validatedInput.CreateMemberAccess
	logEntityProposal, err := createMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "MemberAccess",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityCreate,
			ProposalStatus:   *validatedInput.CreateMemberAccess.ProposalStatus,
			NewObject:        &newObject,
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}

	validatedInput.CreateMemberAccess.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.CreateMemberAccess.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	if *validatedInput.CreateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
		validatedInput.CreateMemberAccess.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
	}
	createdMemberAccess, err := createMmbAccessRefUcase.createMemberAccessRepo.Execute(
		validatedInput.CreateMemberAccess,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}

	return createdMemberAccess, nil
}
