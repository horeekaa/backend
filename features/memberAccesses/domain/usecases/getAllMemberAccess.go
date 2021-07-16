package memberaccesspresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllMemberAccessUsecase struct {
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllMemberAccessRepo     memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository

	getOwnedMemberAccessAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllMemberAccessUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository,
) (memberaccesspresentationusecaseinterfaces.GetAllMemberAccessUsecase, error) {
	return &getAllMemberAccessUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllMemberAccessRepo,
		&model.MemberAccessRefOptionsInput{
			ManageMemberAccesses: &model.ManageMemberAccessesInput{
				MemberAccessReadOwned: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllMmbAccUcase *getAllMemberAccessUsecase) validation(input memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput) (*memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput, error) {
	if &input.Context == nil {
		return &memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllMemberAccessUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllMmbAccUcase *getAllMemberAccessUsecase) Execute(
	input memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput,
) ([]*model.MemberAccess, error) {
	validatedInput, err := getAllMmbAccUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllMmbAccUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMemberAccessUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllMemberAccessUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrganization := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllMmbAccUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
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
			"/getAllMemberAccessUsecase",
			err,
		)
	}

	if accessible, ok := (funk.Get(memberAccess, "Access.ManageMemberAccesses.MemberAccessReadAll")).(bool); !ok || !accessible {
		accessible, ok := (funk.Get(memberAccess, "Access.ManageMemberAccesses.MemberAccessReadOrganizationOwned")).(bool)
		if accessible && ok {
			validatedInput.FilterFields.Organization = &model.AttachOrganizationInput{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			memberAccessRefTypeAccountBasics := model.MemberAccessRefTypeAccountsBasics
			_, err := getAllMmbAccUcase.getAccountMemberAccessRepo.Execute(
				memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
					MemberAccessFilterFields: &model.MemberAccessFilterFields{
						Account:             &model.ObjectIDOnly{ID: &account.ID},
						MemberAccessRefType: &memberAccessRefTypeAccountBasics,
						Access:              getAllMmbAccUcase.getOwnedMemberAccessAccessIdentity,
					},
				},
			)
			if err != nil {
				return nil, horeekaacorefailuretoerror.ConvertFailure(
					"/getAllMemberAccessUsecase",
					err,
				)
			}

			validatedInput.FilterFields.Account = &model.ObjectIDOnly{
				ID: &account.ID,
			}
		}
	}

	memberAccesses, err := getAllMmbAccUcase.getAllMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAllMemberAccessInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMemberAccessUsecase",
			err,
		)
	}

	return memberAccesses, nil
}
