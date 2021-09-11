package moupresentationusecases

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
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositorytypes "github.com/horeekaa/backend/features/mous/domain/repositories/types"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllMouUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllmouRepo              moudomainrepositoryinterfaces.GetAllMouRepository

	getOwnedmouIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllMouUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllmouRepo moudomainrepositoryinterfaces.GetAllMouRepository,
) (moupresentationusecaseinterfaces.GetAllMouUsecase, error) {
	return &getAllMouUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllmouRepo,
		&model.MemberAccessRefOptionsInput{
			MouAccesses: &model.MouAccessesInput{
				MouReadOwned: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllMouUcase *getAllMouUsecase) validation(input moupresentationusecasetypes.GetAllMouUsecaseInput) (*moupresentationusecasetypes.GetAllMouUsecaseInput, error) {
	if &input.Context == nil {
		return &moupresentationusecasetypes.GetAllMouUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllMouUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllMouUcase *getAllMouUsecase) Execute(
	input moupresentationusecasetypes.GetAllMouUsecaseInput,
) ([]*model.Mou, error) {
	validatedInput, err := getAllMouUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllMouUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMouUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllMouUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllMouUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
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
			"/getAllMouUsecase",
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.MouAccesses.MouReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.MouAccesses.MouReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.SecondParty = &model.PartyInput{
				Organization: &model.ObjectIDOnly{
					ID: &memberAccess.Organization.ID,
				},
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getAllMouUsecase",
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					"/getAllMouUsecase",
					nil,
				),
			)
		}
	}

	mous, err := getAllMouUcase.getAllmouRepo.Execute(
		moudomainrepositorytypes.GetAllMouInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMouUsecase",
			err,
		)
	}

	return mous, nil
}
