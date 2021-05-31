package organizationpresentationusecases

import (
	"fmt"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
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
	getPersonDataFromAccountRepo     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	updateOrganizationRepo           organizationdomainrepositoryinterfaces.UpdateOrganizationRepository
	getOrganizationRepo              organizationdomainrepositoryinterfaces.GetOrganizationRepository
	getAllMemberAccessRepo           memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository
	updateMemberAccessRepo           memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository
	logEntityProposalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository
	logEntityApprovalActivityRepo    loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository
	updateOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
	updateOrganizationRepo organizationdomainrepositoryinterfaces.UpdateOrganizationRepository,
	getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
	getAllMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository,
	updateMemberAccessRepo memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository,
	logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
	logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
) (organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase, error) {
	return &updateOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getPersonDataFromAccountRepo,
		updateOrganizationRepo,
		getOrganizationRepo,
		getAllMemberAccessRepo,
		updateMemberAccessRepo,
		logEntityProposalActivityRepo,
		logEntityApprovalActivityRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateMmbAccessRefUcase *updateOrganizationUsecase) validation(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput, error) {
	if &input.Context == nil {
		return organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateOrganizationUsecase",
				nil,
			)
	}
	input.UpdateOrganization.ApprovingAccount = nil
	input.UpdateOrganization.SubmittingAccount = nil

	return input, nil
}

func (updateMmbAccessRefUcase *updateOrganizationUsecase) Execute(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (*model.Organization, error) {
	validatedInput, err := updateMmbAccessRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateMmbAccessRefUcase.getAccountFromAuthDataRepo.Execute(
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

	personChannel := make(chan *model.Person)
	errChannel := make(chan error)
	go func() {
		person, err := updateMmbAccessRefUcase.getPersonDataFromAccountRepo.Execute(account)
		if err != nil {
			errChannel <- err
		}
		personChannel <- person
	}()

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateMmbAccessRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganization,
				Access:              updateMmbAccessRefUcase.updateOrganizationAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
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
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	existingOrg, err := updateMmbAccessRefUcase.getOrganizationRepo.Execute(
		&model.OrganizationFilterFields{
			ID: &validatedInput.UpdateOrganization.ID,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	// if user is only going to approve proposal
	if validatedInput.UpdateOrganization.ProposalStatus != nil {
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

		logApprovalActivity, err := updateMmbAccessRefUcase.logEntityApprovalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityApprovalActivityInput{
				PreviousLog:      existingOrg.CorrespondingLog,
				ApprovingAccount: account,
				ApproverInitial:  accountInitials,
				ApprovalStatus:   *validatedInput.UpdateOrganization.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		validatedInput.UpdateOrganization.ApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		validatedInput.UpdateOrganization.CorrespondingLog = &model.ObjectIDOnly{ID: &logApprovalActivity.ID}
		updateOrganizationOutput, err := updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
			validatedInput.UpdateOrganization,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}
		_, err = updateMmbAccessRefUcase.updateCorrespondingMemberAccess(
			existingOrg,
			updateOrganizationOutput.UpdatedOrganization,
			account,
			accountInitials,
		)
		if err != nil {
			return nil, err
		}

		return updateOrganizationOutput.UpdatedOrganization, nil
	}

	validatedInput.UpdateOrganization.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.OrganizationAccesses.OrganizationApproval != nil {
		if *accMemberAccess.Access.OrganizationAccesses.OrganizationApproval {
			validatedInput.UpdateOrganization.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	var newObject interface{} = *validatedInput.UpdateOrganization
	var existingObject interface{} = *existingOrg
	logEntityProposal, err := updateMmbAccessRefUcase.logEntityProposalActivityRepo.Execute(
		loggingdomainrepositorytypes.LogEntityProposalActivityInput{
			CollectionName:   "Organization",
			CreatedByAccount: account,
			Activity:         model.LoggedActivityUpdate,
			ProposalStatus:   *validatedInput.UpdateOrganization.ProposalStatus,
			NewObject:        &newObject,
			ExistingObject:   &existingObject,
			ExistingObjectID: func(t string) *string { return &t }(existingOrg.ID.Hex()),
			CreatorInitial:   accountInitials,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	validatedInput.UpdateOrganization.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	validatedInput.UpdateOrganization.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}
	updateOrganizationOutput, err := updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
		validatedInput.UpdateOrganization,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	// user is going to update and directly has permission to approve the update
	if *validatedInput.UpdateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
		updateOrganizationOutput, err = updateMmbAccessRefUcase.updateOrganizationRepo.RunTransaction(
			&model.UpdateOrganization{
				ID:               updateOrganizationOutput.UpdatedOrganization.ID,
				ApprovingAccount: &model.ObjectIDOnly{ID: &account.ID},
				ProposalStatus:   validatedInput.UpdateOrganization.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}
		_, err := updateMmbAccessRefUcase.updateCorrespondingMemberAccess(
			existingOrg,
			updateOrganizationOutput.UpdatedOrganization,
			account,
			accountInitials,
		)
		if err != nil {
			return nil, err
		}

		return updateOrganizationOutput.UpdatedOrganization, nil
	}

	return updateOrganizationOutput.UpdatedOrganization, nil
}

func (updateOrgUcase *updateOrganizationUsecase) updateCorrespondingMemberAccess(
	existingOrg *model.Organization,
	updatedOrg *model.Organization,
	account *model.Account,
	accountInitials string,
) (*bool, error) {
	memberAccessesToUpdate, err := updateOrgUcase.getAllMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAllMemberAccessInput{
			FilterFields: &model.MemberAccessFilterFields{
				Organization: &model.AttachOrganizationInput{
					ID: &existingOrg.ID,
				},
				Status: func(s model.MemberAccessStatus) *model.MemberAccessStatus {
					return &s
				}(model.MemberAccessStatusActive),
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateOrganizationUsecase",
			err,
		)
	}

	for _, memberAccess := range memberAccessesToUpdate {
		if memberAccess.ProposalStatus == model.EntityProposalStatusReplaced {
			continue
		}

		updateMemberAccessData := &model.UpdateMemberAccess{
			ID: memberAccess.ID,
			Organization: &model.AttachOrganizationInput{
				ID:   &updatedOrg.ID,
				Type: &updatedOrg.Type,
			},
			SubmittingAccount: &model.ObjectIDOnly{ID: &account.ID},
			ProposalStatus: func(ep model.EntityProposalStatus) *model.EntityProposalStatus {
				return &ep
			}(model.EntityProposalStatusApproved),
		}

		var newObject interface{} = *updateMemberAccessData
		var existingObject interface{} = *memberAccess
		logEntityProposal, err := updateOrgUcase.logEntityProposalActivityRepo.Execute(
			loggingdomainrepositorytypes.LogEntityProposalActivityInput{
				CollectionName:   "MemberAccess",
				CreatedByAccount: account,
				Activity:         model.LoggedActivityUpdate,
				ProposalStatus:   *updateMemberAccessData.ProposalStatus,
				NewObject:        &newObject,
				ExistingObject:   &existingObject,
				ExistingObjectID: func(t string) *string { return &t }(memberAccess.ID.Hex()),
				CreatorInitial:   accountInitials,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}
		updateMemberAccessData.CorrespondingLog = &model.ObjectIDOnly{ID: &logEntityProposal.ID}

		updatedMemberAccess, err := updateOrgUcase.updateMemberAccessRepo.RunTransaction(
			updateMemberAccessData,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}

		_, err = updateOrgUcase.updateMemberAccessRepo.RunTransaction(
			&model.UpdateMemberAccess{
				ID:               updatedMemberAccess.UpdatedMemberAccess.ID,
				ApprovingAccount: &model.ObjectIDOnly{ID: &account.ID},
				ProposalStatus:   &updatedMemberAccess.UpdatedMemberAccess.ProposalStatus,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/updateOrganizationUsecase",
				err,
			)
		}
	}
	return func(b bool) *bool { return &b }(true), nil
}
