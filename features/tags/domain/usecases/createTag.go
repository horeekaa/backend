package tagpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createTagUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createTagRepo              tagdomainrepositoryinterfaces.CreateTagRepository
	createTagAccessIdentity    *model.MemberAccessRefOptionsInput
}

func NewCreateTagUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createTagRepo tagdomainrepositoryinterfaces.CreateTagRepository,
) (tagpresentationusecaseinterfaces.CreateTagUsecase, error) {
	return &createTagUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createTagRepo,
		&model.MemberAccessRefOptionsInput{
			TagAccesses: &model.TagAccessesInput{
				TagCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createTagUcase *createTagUsecase) validation(input tagpresentationusecasetypes.CreateTagUsecaseInput) (tagpresentationusecasetypes.CreateTagUsecaseInput, error) {
	if &input.Context == nil {
		return tagpresentationusecasetypes.CreateTagUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createTagUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateTag.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createTagUcase *createTagUsecase) Execute(input tagpresentationusecasetypes.CreateTagUsecaseInput) (*model.Tag, error) {
	validatedInput, err := createTagUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createTagUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createTagUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createTagUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createTagUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createTagUcase.createTagAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createTagUsecase",
			err,
		)
	}
	if accMemberAccess.Access.TagAccesses.TagApproval != nil {
		if *accMemberAccess.Access.TagAccesses.TagApproval {
			validatedInput.CreateTag.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	tagToCreate := &model.InternalCreateTag{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateTag)
	json.Unmarshal(jsonTemp, tagToCreate)

	for i, descriptivePhoto := range validatedInput.CreateTag.Photos {
		if descriptivePhoto.Photo != nil {
			tagToCreate.Photos[i].Photo.File = descriptivePhoto.Photo.File
		}
	}

	tagToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdTag, err := createTagUcase.createTagRepo.RunTransaction(
		tagToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createTagUsecase",
			err,
		)
	}

	return createdTag, nil
}
