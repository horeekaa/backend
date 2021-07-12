package organizationpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateOrganizationUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateOrganizationRepo    organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository
	approveUpdateOrganizationRepo    organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository
	updateOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository,
	approveUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository,
) (organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase, error) {
	return &updateOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateOrganizationRepo,
		approveUpdateOrganizationRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateOrganizationUcase *updateOrganizationUsecase) validation(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput, error) {
	if &input.Context == nil {
		return organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateOrganizationUsecase",
				nil,
			)
	}

	return input, nil
}

func (updateOrganizationUcase *updateOrganizationUsecase) Execute(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := updateOrganizationUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateOrganizationUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updateOrganizationUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateOrganizationUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Access:              updateOrganizationUcase.updateOrganizationAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	organizationToUpdate := &model.InternalUpdateOrganization{
		ID: validatedInput.UpdateOrganization.ID,
	}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateOrganization)
	json.Unmarshal(jsonTemp, organizationToUpdate)

	// if user is only going to approve proposal
	if organizationToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.OrganizationAccesses.OrganizationApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateOrganizationUsecase",
				nil,
			)
		}
		if !*accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				403,
				"/updateOrganizationUsecase",
				nil,
			)
		}

		organizationToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateOrganizationOutput, err := updateOrganizationUcase.approveUpdateOrganizationRepo.RunTransaction(
			organizationToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		return updateOrganizationOutput, nil
	}

	organizationToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.OrganizationAccesses.OrganizationApproval != nil {
		if *accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
			organizationToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	organizationToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateOrganizationOutput, err := updateOrganizationUcase.proposeUpdateOrganizationRepo.RunTransaction(
		organizationToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	return updateOrganizationOutput, nil
}
