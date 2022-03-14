package taggingpresentationusecases

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
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type bulkUpdateTaggingUsecase struct {
	getAccountFromAuthDataRepo   accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo   memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeBulkUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingRepository
	approveBulkUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingRepository
	updateTaggingAccessIdentity  *model.MemberAccessRefOptionsInput
	pathIdentity                 string
}

func NewBulkUpdateTaggingUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	bulkProposeUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingRepository,
	bulkApproveUpdateTaggingRepo taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingRepository,
) (taggingpresentationusecaseinterfaces.BulkUpdateTaggingUsecase, error) {
	return &bulkUpdateTaggingUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		bulkProposeUpdateTaggingRepo,
		bulkApproveUpdateTaggingRepo,
		&model.MemberAccessRefOptionsInput{
			BulkTaggingAccesses: &model.BulkTaggingAccessesInput{
				BulkTaggingUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"BulkUpdateTaggingUsecase",
	}, nil
}

func (updateTaggingUcase *bulkUpdateTaggingUsecase) validation(input taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput) (taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput, error) {
	if &input.Context == nil {
		return taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateTaggingUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateTaggingUcase *bulkUpdateTaggingUsecase) Execute(input taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput) ([]*model.Tagging, error) {
	validatedInput, err := updateTaggingUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateTaggingUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTaggingUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateTaggingUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateTaggingUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateTaggingUcase.updateTaggingAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTaggingUcase.pathIdentity,
			err,
		)
	}

	taggingToUpdate := &model.InternalBulkUpdateTagging{
		IDs: validatedInput.BulkUpdateTagging.IDs,
	}
	jsonTemp, _ := json.Marshal(validatedInput.BulkUpdateTagging)
	json.Unmarshal(jsonTemp, taggingToUpdate)

	// if user is only going to approve proposal
	if taggingToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateTaggingUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateTaggingUcase.pathIdentity,
				nil,
			)
		}

		taggingToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updatedTaggings, err := updateTaggingUcase.approveBulkUpdateTaggingRepo.RunTransaction(
			taggingToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateTaggingUcase.pathIdentity,
				err,
			)
		}

		return updatedTaggings, nil
	}

	taggingToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval != nil {
		if *accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval {
			taggingToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	taggingToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updatedTaggings, err := updateTaggingUcase.proposeBulkUpdateTaggingRepo.RunTransaction(
		taggingToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTaggingUcase.pathIdentity,
			err,
		)
	}

	return updatedTaggings, nil
}
