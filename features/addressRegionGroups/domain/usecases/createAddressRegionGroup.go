package addressregiongrouppresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type createAddressRegionGroupUsecase struct {
	getAccountFromAuthDataRepo             accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo             memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createAddressRegionGroupRepo           addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupRepository
	createAddressRegionGroupAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateAddressRegionGroupUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createAddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupRepository,
) (addressregiongrouppresentationusecaseinterfaces.CreateAddressRegionGroupUsecase, error) {
	return &createAddressRegionGroupUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createAddressRegionGroupRepo,
		&model.MemberAccessRefOptionsInput{
			AddressRegionGroupAccesses: &model.AddressRegionGroupAccessesInput{
				AddressRegionGroupCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createAddressRegionGroupUcase *createAddressRegionGroupUsecase) validation(input addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput) (addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput, error) {
	if &input.Context == nil {
		return addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createAddressRegionGroupUsecase",
				nil,
			)
	}
	return input, nil
}

func (createAddressRegionGroupUcase *createAddressRegionGroupUsecase) Execute(input addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput) (*model.AddressRegionGroup, error) {
	validatedInput, err := createAddressRegionGroupUcase.validation(input)
	if err != nil {
		return nil, err
	}
	addressRegionGroupToCreate := &model.InternalCreateAddressRegionGroup{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateAddressRegionGroup)
	json.Unmarshal(jsonTemp, addressRegionGroupToCreate)

	account, err := createAddressRegionGroupUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createAddressRegionGroupUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createAddressRegionGroupUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createAddressRegionGroupUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createAddressRegionGroupUcase.createAddressRegionGroupAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createAddressRegionGroupUsecase",
			err,
		)
	}

	addressRegionGroupToCreate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval != nil {
		if *accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval {
			addressRegionGroupToCreate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	addressRegionGroupToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdAddressRegionGroup, err := createAddressRegionGroupUcase.createAddressRegionGroupRepo.RunTransaction(
		addressRegionGroupToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createAddressRegionGroupUsecase",
			err,
		)
	}

	return createdAddressRegionGroup, nil
}
