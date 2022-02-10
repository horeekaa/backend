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
	"github.com/thoas/go-funk"
)

type updateOrganizationUsecase struct {
	getAccountFromAuthDataRepo    accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo    memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository
	approveUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository
	createMemberAccessRepo        memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository
	getOrganizationRepo           organizationdomainrepositoryinterfaces.GetOrganizationRepository

	updateOwnedOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository,
	approveUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository,
	createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
	getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
) (organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase, error) {
	return &updateOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateOrganizationRepo,
		approveUpdateOrganizationRepo,
		createMemberAccessRepo,
		getOrganizationRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationUpdateOwned: func(b bool) *bool { return &b }(true),
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
	organizationToUpdate := &model.InternalUpdateOrganization{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateOrganization)
	json.Unmarshal(jsonTemp, organizationToUpdate)

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
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Status: func(s model.MemberAccessStatus) *model.MemberAccessStatus {
					return &s
				}(model.MemberAccessStatusActive),
				ProposalStatus: func(e model.EntityProposalStatus) *model.EntityProposalStatus {
					return &e
				}(model.EntityProposalStatusApproved),
				InvitationAccepted: func(b bool) *bool {
					return &b
				}(true),
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	existingOrganization, err := updateOrganizationUcase.getOrganizationRepo.Execute(
		&model.OrganizationFilterFields{
			ID: &organizationToUpdate.ID,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}
	// if update across organizations is not allowed, check access for update owned organization
	if accessible := funk.GetOrElse(
		funk.Get(accMemberAccess, "Access.OrganizationAccesses.OrganizationUpdate"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(accMemberAccess, "Access.OrganizationAccesses.OrganizationUpdateOwned"), false,
		).(bool); accessible {
			if accMemberAccess.Organization.ID != organizationToUpdate.ID {
				return nil, horeekaacoreerror.NewErrorObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					403,
					"/updateOrganizationUsecase",
					nil,
				)
			}
		} else {
			memberAccessRefTypeAccountBasics := model.MemberAccessRefTypeAccountsBasics
			accMemberAccess, err = updateOrganizationUcase.getAccountMemberAccessRepo.Execute(
				memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
					MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
						Account:             &model.ObjectIDOnly{ID: &account.ID},
						Access:              updateOrganizationUcase.updateOwnedOrganizationAccessIdentity,
						MemberAccessRefType: &memberAccessRefTypeAccountBasics,
					},
				},
			)
			if err != nil {
				return nil, horeekaacorefailuretoerror.ConvertFailure(
					"/updateOrganizationUsecase",
					err,
				)
			}
		}
	}

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

		if existingOrganization.ProposalStatus != model.EntityProposalStatusApproved &&
			updateOrganizationOutput.ProposalStatus == model.EntityProposalStatusApproved {
			_, err := updateOrganizationUcase.createMemberAccessRepo.RunTransaction(
				&model.InternalCreateMemberAccess{
					Account: &model.ObjectIDOnly{
						ID: &existingOrganization.SubmittingAccount.ID,
					},
					Organization: &model.InternalUpdateOrganization{
						ID: updateOrganizationOutput.ID,
					},
					OrganizationMembershipRole: func(r model.OrganizationMembershipRole) *model.OrganizationMembershipRole {
						return &r
					}(model.OrganizationMembershipRoleOwner),
					MemberAccessRefType: model.MemberAccessRefTypeOrganizationsBased,
					SubmittingAccount: &model.ObjectIDOnly{
						ID: &account.ID,
					},
					ProposalStatus: func(e model.EntityProposalStatus) *model.EntityProposalStatus {
						return &e
					}(model.EntityProposalStatusApproved),
					InvitationAccepted: func(b bool) *bool {
						return &b
					}(true),
				},
			)
			if err != nil {
				return nil, horeekaacorefailuretoerror.ConvertFailure(
					"/updateOrganizationUsecase",
					err,
				)
			}
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
