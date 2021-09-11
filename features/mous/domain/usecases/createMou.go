package moupresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createMouUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createMouRepo              moudomainrepositoryinterfaces.CreateMouRepository
	createMouAccessIdentity    *model.MemberAccessRefOptionsInput
}

func NewCreateMouUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createMouRepo moudomainrepositoryinterfaces.CreateMouRepository,
) (moupresentationusecaseinterfaces.CreateMouUsecase, error) {
	return &createMouUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createMouRepo,
		&model.MemberAccessRefOptionsInput{
			MouAccesses: &model.MouAccessesInput{
				MouCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createMouUcase *createMouUsecase) validation(input moupresentationusecasetypes.CreateMouUsecaseInput) (moupresentationusecasetypes.CreateMouUsecaseInput, error) {
	if &input.Context == nil {
		return moupresentationusecasetypes.CreateMouUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createMouUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateMou.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createMouUcase *createMouUsecase) Execute(input moupresentationusecasetypes.CreateMouUsecaseInput) (*model.Mou, error) {
	validatedInput, err := createMouUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createMouUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMouUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createMouUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createMouUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createMouUcase.createMouAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMouUsecase",
			err,
		)
	}
	if accMemberAccess.Access.MouAccesses.MouApproval != nil {
		if *accMemberAccess.Access.MouAccesses.MouApproval {
			validatedInput.CreateMou.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	mouToCreate := &model.InternalCreateMou{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateMou)
	json.Unmarshal(jsonTemp, mouToCreate)

	mouToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdMou, err := createMouUcase.createMouRepo.RunTransaction(
		mouToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createMouUsecase",
			err,
		)
	}

	return createdMou, nil
}
