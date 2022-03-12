package tagpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagdomainrepositorytypes "github.com/horeekaa/backend/features/tags/domain/repositories/types"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllTagUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllTagRepo              tagdomainrepositoryinterfaces.GetAllTagRepository
	getAllTagAccessIdentity    *model.MemberAccessRefOptionsInput
	pathIdentity               string
}

func NewGetAllTagUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllTagRepo tagdomainrepositoryinterfaces.GetAllTagRepository,
) (tagpresentationusecaseinterfaces.GetAllTagUsecase, error) {
	return &getAllTagUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllTagRepo,
		&model.MemberAccessRefOptionsInput{
			TagAccesses: &model.TagAccessesInput{
				TagReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllTagUsecase",
	}, nil
}

func (getAllTagUcase *getAllTagUsecase) validation(input tagpresentationusecasetypes.GetAllTagUsecaseInput) (*tagpresentationusecasetypes.GetAllTagUsecaseInput, error) {
	if &input.Context == nil {
		return &tagpresentationusecasetypes.GetAllTagUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllTagUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllTagUcase *getAllTagUsecase) Execute(
	input tagpresentationusecasetypes.GetAllTagUsecaseInput,
) ([]*model.Tag, error) {
	validatedInput, err := getAllTagUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllTagUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTagUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllTagUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllTagUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllTagUcase.getAllTagAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTagUcase.pathIdentity,
			err,
		)
	}

	tags, err := getAllTagUcase.getAllTagRepo.Execute(
		tagdomainrepositorytypes.GetAllTagInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTagUcase.pathIdentity,
			err,
		)
	}

	return tags, nil
}
