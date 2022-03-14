package addressregiongrouppresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateAddressRegionGroupUsecase struct {
	getAccountFromAuthDataRepo             accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo             memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateAddressRegionGroupRepo    addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupRepository
	approveUpdateAddressRegionGroupRepo    addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupRepository
	updateAddressRegionGroupAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                           string
}

func NewUpdateAddressRegionGroupUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateAddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupRepository,
	approveUpdateAddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupRepository,
) (addressregiongrouppresentationusecaseinterfaces.UpdateAddressRegionGroupUsecase, error) {
	return &updateAddressRegionGroupUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateAddressRegionGroupRepo,
		approveUpdateAddressRegionGroupRepo,
		&model.MemberAccessRefOptionsInput{
			AddressRegionGroupAccesses: &model.AddressRegionGroupAccessesInput{
				AddressRegionGroupUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateAddressRegionGroupUsecase",
	}, nil
}

func (updateAddressRegionGroupUcase *updateAddressRegionGroupUsecase) validation(input addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput) (addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput, error) {
	if &input.Context == nil {
		return addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateAddressRegionGroupUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateAddressRegionGroupUcase *updateAddressRegionGroupUsecase) Execute(input addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput) (*model.AddressRegionGroup, error) {
	validatedInput, err := updateAddressRegionGroupUcase.validation(input)
	if err != nil {
		return nil, err
	}
	addressRegionGroupToUpdate := &model.InternalUpdateAddressRegionGroup{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateAddressRegionGroup)
	json.Unmarshal(jsonTemp, addressRegionGroupToUpdate)

	account, err := updateAddressRegionGroupUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateAddressRegionGroupUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateAddressRegionGroupUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateAddressRegionGroupUcase.updateAddressRegionGroupAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}

	// if user is only going to approve proposal
	if addressRegionGroupToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateAddressRegionGroupUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateAddressRegionGroupUcase.pathIdentity,
				nil,
			)
		}

		addressRegionGroupToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateAddressRegionGroupOutput, err := updateAddressRegionGroupUcase.approveUpdateAddressRegionGroupRepo.RunTransaction(
			addressRegionGroupToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateAddressRegionGroupUcase.pathIdentity,
				err,
			)
		}

		return updateAddressRegionGroupOutput, nil
	}

	addressRegionGroupToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval != nil {
		if *accMemberAccess.Access.AddressRegionGroupAccesses.AddressRegionGroupApproval {
			addressRegionGroupToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	addressRegionGroupToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateAddressRegionGroupOutput, err := updateAddressRegionGroupUcase.proposeUpdateAddressRegionGroupRepo.RunTransaction(
		addressRegionGroupToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}

	return updateAddressRegionGroupOutput, nil
}
