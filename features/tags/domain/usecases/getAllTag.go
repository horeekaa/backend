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
	}, nil
}

func (getAllProdUcase *getAllTagUsecase) validation(input tagpresentationusecasetypes.GetAllTagUsecaseInput) (*tagpresentationusecasetypes.GetAllTagUsecaseInput, error) {
	if &input.Context == nil {
		return &tagpresentationusecasetypes.GetAllTagUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllTagUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllProdUcase *getAllTagUsecase) Execute(
	input tagpresentationusecasetypes.GetAllTagUsecaseInput,
) ([]*model.Tag, error) {
	validatedInput, err := getAllProdUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllProdUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllTagUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllTagUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllProdUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllProdUcase.getAllTagAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllTagUsecase",
			err,
		)
	}

	tags, err := getAllProdUcase.getAllTagRepo.Execute(
		tagdomainrepositorytypes.GetAllTagInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllTagUsecase",
			err,
		)
	}

	return tags, nil
}
