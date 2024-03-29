package memberaccessrefpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefUsecase struct {
	getAccountFromAuthDataRepo          accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo          memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository
	createMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                        string
}

func NewCreateMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository,
) (memberaccessrefpresentationusecaseinterfaces.CreateMemberAccessRefUsecase, error) {
	return &createMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createMemberAccessRefRepo,
		&model.MemberAccessRefOptionsInput{
			MemberAccessRefAccesses: &model.MemberAccessRefAccessesInput{
				MemberAccessRefCreate: func(b bool) *bool { return &b }(true),
			},
		},
		"CreateMemberAccessRefUsecase",
	}, nil
}

func (createMmbAccessRefUcase *createMemberAccessRefUsecase) validation(input memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput) (memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput, error) {
	if &input.Context == nil {
		return memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				createMmbAccessRefUcase.pathIdentity,
				nil,
			)
	}
	return input, nil
}

func (createMmbAccessRefUcase *createMemberAccessRefUsecase) Execute(input memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput) (*model.MemberAccessRef, error) {
	validatedInput, err := createMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}
	memberAccessRefToCreate := &model.InternalCreateMemberAccessRef{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateMemberAccessRef)
	json.Unmarshal(jsonTemp, memberAccessRefToCreate)

	account, err := createMmbAccessRefUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			createMmbAccessRefUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			createMmbAccessRefUcase.pathIdentity,
			nil,
		)
	}

	accMemberAccess, err := createMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account: &model.ObjectIDOnly{ID: &account.ID},
				Access:  createMmbAccessRefUcase.createMemberAccessRefAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			createMmbAccessRefUcase.pathIdentity,
			err,
		)
	}

	memberAccessRefToCreate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval != nil {
		if *accMemberAccess.Access.MemberAccessRefAccesses.MemberAccessRefApproval {
			memberAccessRefToCreate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	memberAccessRefToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdMemberAccessRef, err := createMmbAccessRefUcase.createMemberAccessRefRepo.RunTransaction(
		memberAccessRefToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			createMmbAccessRefUcase.pathIdentity,
			err,
		)
	}

	return createdMemberAccessRef, nil
}
