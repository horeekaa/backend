package tagpresentationusecases

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
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateTagUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateTagRepo       tagdomainrepositoryinterfaces.ProposeUpdateTagRepository
	approveUpdateTagRepo       tagdomainrepositoryinterfaces.ApproveUpdateTagRepository
	updateTagAccessIdentity    *model.MemberAccessRefOptionsInput
	pathIdentity               string
}

func NewUpdateTagUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateTagRepo tagdomainrepositoryinterfaces.ProposeUpdateTagRepository,
	approveUpdateTagRepo tagdomainrepositoryinterfaces.ApproveUpdateTagRepository,
) (tagpresentationusecaseinterfaces.UpdateTagUsecase, error) {
	return &updateTagUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateTagRepo,
		approveUpdateTagRepo,
		&model.MemberAccessRefOptionsInput{
			TagAccesses: &model.TagAccessesInput{
				TagUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateTagUsecase",
	}, nil
}

func (updateTagUcase *updateTagUsecase) validation(input tagpresentationusecasetypes.UpdateTagUsecaseInput) (tagpresentationusecasetypes.UpdateTagUsecaseInput, error) {
	if &input.Context == nil {
		return tagpresentationusecasetypes.UpdateTagUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateTagUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateTagUcase *updateTagUsecase) Execute(input tagpresentationusecasetypes.UpdateTagUsecaseInput) (*model.Tag, error) {
	validatedInput, err := updateTagUcase.validation(input)
	if err != nil {
		return nil, err
	}
	tagToUpdate := &model.InternalUpdateTag{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateTag)
	json.Unmarshal(jsonTemp, tagToUpdate)

	account, err := updateTagUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTagUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateTagUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateTagUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateTagUcase.updateTagAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTagUcase.pathIdentity,
			err,
		)
	}

	for i, descriptivePhoto := range validatedInput.UpdateTag.Photos {
		if descriptivePhoto.Photo != nil {
			tagToUpdate.Photos[i].Photo.File = descriptivePhoto.Photo.File
		}
	}

	// if user is only going to approve proposal
	if tagToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.TagAccesses.TagApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateTagUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.TagAccesses.TagApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateTagUcase.pathIdentity,
				nil,
			)
		}

		tagToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateTagOutput, err := updateTagUcase.approveUpdateTagRepo.RunTransaction(
			tagToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateTagUcase.pathIdentity,
				err,
			)
		}

		return updateTagOutput, nil
	}

	tagToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.TagAccesses.TagApproval != nil {
		if *accMemberAccess.Access.TagAccesses.TagApproval {
			tagToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	tagToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateTagOutput, err := updateTagUcase.proposeUpdateTagRepo.RunTransaction(
		tagToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateTagUcase.pathIdentity,
			err,
		)
	}

	return updateTagOutput, nil
}
