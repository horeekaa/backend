package memberaccesspresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createMemberAccessRepo           memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository
	createMemberAccessAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateMemberAccessUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
) (memberaccesspresentationusecaseinterfaces.CreateMemberAccessUsecase, error) {
	return &createMemberAccessUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createMemberAccessRepo,
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMmbAccessUcase *createMemberAccessUsecase) validation(input memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput) (memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput, error) {
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
	input.CreateMemberAccess.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createMmbAccessUcase *createMemberAccessUsecase) Execute(input memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput) (*model.MemberAccess, error) {
	validatedInput, err := createMmbAccessUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMmbAccessUcase.getAccountFromAuthDataRepo.Execute(
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

	duplicateMemberAccess, err := createMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account: &model.ObjectIDOnly{ID: validatedInput.CreateMemberAccess.Account.ID},
				Status: func(m model.MemberAccessStatus) *model.MemberAccessStatus {
					return &m
				}(model.MemberAccessStatusActive),
				ProposalStatus: func(m model.EntityProposalStatus) *model.EntityProposalStatus {
					return &m
				}(model.EntityProposalStatusApproved),
				MemberAccessRefType: &validatedInput.CreateMemberAccess.MemberAccessRefType,
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
			409,
			"/createMemberAccessUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createMmbAccessUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createMmbAccessUcase.createMemberAccessAccessIdentity,
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

	memberAccessToCreate := &model.InternalCreateMemberAccess{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateMemberAccess)
	json.Unmarshal(jsonTemp, memberAccessToCreate)

	memberAccessToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdMemberAccess, err := createMmbAccessUcase.createMemberAccessRepo.RunTransaction(
		memberAccessToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMemberAccessUsecase",
			err,
		)
	}

	return createdMemberAccess, nil
}
