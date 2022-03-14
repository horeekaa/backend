package taggingpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingdomainrepositorytypes "github.com/horeekaa/backend/features/taggings/domain/repositories/types"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"

	"github.com/horeekaa/backend/model"
)

type getAllTaggingUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllTaggingRepo           taggingdomainrepositoryinterfaces.GetAllTaggingRepository
	getAllTaggingAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewGetAllTaggingUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllTaggingRepo taggingdomainrepositoryinterfaces.GetAllTaggingRepository,
) (taggingpresentationusecaseinterfaces.GetAllTaggingUsecase, error) {
	return &getAllTaggingUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllTaggingRepo,
		&model.MemberAccessRefOptionsInput{
			BulkTaggingAccesses: &model.BulkTaggingAccessesInput{
				BulkTaggingReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllTaggingUsecase",
	}, nil
}

func (getAllTaggingUcase *getAllTaggingUsecase) validation(input taggingpresentationusecasetypes.GetAllTaggingUsecaseInput) (*taggingpresentationusecasetypes.GetAllTaggingUsecaseInput, error) {
	if &input.Context == nil {
		return &taggingpresentationusecasetypes.GetAllTaggingUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllTaggingUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllTaggingUcase *getAllTaggingUsecase) Execute(
	input taggingpresentationusecasetypes.GetAllTaggingUsecaseInput,
) ([]*model.Tagging, error) {
	validatedInput, err := getAllTaggingUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllTaggingUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTaggingUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllTaggingUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllTaggingUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllTaggingUcase.getAllTaggingAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTaggingUcase.pathIdentity,
			err,
		)
	}

	taggings, err := getAllTaggingUcase.getAllTaggingRepo.Execute(
		taggingdomainrepositorytypes.GetAllTaggingInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllTaggingUcase.pathIdentity,
			err,
		)
	}

	return taggings, nil
}
