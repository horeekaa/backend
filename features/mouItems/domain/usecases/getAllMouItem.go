package mouitempresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositorytypes "github.com/horeekaa/backend/features/mouItems/domain/repositories/types"
	mouitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/mouItems/presentation/usecases"
	mouitempresentationusecasetypes "github.com/horeekaa/backend/features/mouItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllMouItemUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllMouItemRepo          mouitemdomainrepositoryinterfaces.GetAllMouItemRepository

	getOwnedMouItemIdentity *model.MemberAccessRefOptionsInput
	pathIdentity            string
}

func NewGetAllMouItemUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllMouItemRepo mouitemdomainrepositoryinterfaces.GetAllMouItemRepository,
) (mouitempresentationusecaseinterfaces.GetAllMouItemUsecase, error) {
	return &getAllMouItemUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllMouItemRepo,
		&model.MemberAccessRefOptionsInput{
			MouAccesses: &model.MouAccessesInput{
				MouReadOwned: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllMouItemUsecase",
	}, nil
}

func (getAllMouItemUcase *getAllMouItemUsecase) validation(input mouitempresentationusecasetypes.GetAllMouItemUsecaseInput) (*mouitempresentationusecasetypes.GetAllMouItemUsecaseInput, error) {
	if &input.Context == nil {
		return &mouitempresentationusecasetypes.GetAllMouItemUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllMouItemUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllMouItemUcase *getAllMouItemUsecase) Execute(
	input mouitempresentationusecasetypes.GetAllMouItemUsecaseInput,
) ([]*model.MouItem, error) {
	validatedInput, err := getAllMouItemUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllMouItemUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllMouItemUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllMouItemUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllMouItemUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
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
			getAllMouItemUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.MouAccesses.MouReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.MouAccesses.MouReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.Organization = &model.OrganizationForMouItemFilterFields{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAllMouItemUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					getAllMouItemUcase.pathIdentity,
					nil,
				),
			)
		}
	}

	mouItems, err := getAllMouItemUcase.getAllMouItemRepo.Execute(
		mouitemdomainrepositorytypes.GetAllMouItemInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllMouItemUcase.pathIdentity,
			err,
		)
	}

	return mouItems, nil
}
