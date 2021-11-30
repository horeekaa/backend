package supplyorderpresentationusecases

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
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/types"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllSupplyOrderUsecase struct {
	getAccountFromAuthDataRepo      accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo      memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllSupplyOrderRepo           supplyorderdomainrepositoryinterfaces.GetAllSupplyOrderRepository
	getAllSupplyOrderAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllSupplyOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllsupplyOrderRepo supplyorderdomainrepositoryinterfaces.GetAllSupplyOrderRepository,
) (supplyorderpresentationusecaseinterfaces.GetAllSupplyOrderUsecase, error) {
	return &getAllSupplyOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllsupplyOrderRepo,
		&model.MemberAccessRefOptionsInput{
			SupplyOrderAccesses: &model.SupplyOrderAccessesInput{
				SupplyOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllSOUcase *getAllSupplyOrderUsecase) validation(input supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput) (*supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput, error) {
	if &input.Context == nil {
		return &supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllSupplyOrderUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllSOUcase *getAllSupplyOrderUsecase) Execute(
	input supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput,
) ([]*model.SupplyOrder, error) {
	validatedInput, err := getAllSOUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllSOUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllSupplyOrderUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllSupplyOrderUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllSOUcase.getAccountMemberAccessRepo.Execute(
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
			"/getAllSupplyOrderUsecase",
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.SupplyOrderAccesses.SupplyOrderReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.SupplyOrderAccesses.SupplyOrderReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.Organization = &model.OrganizationForSupplyOrderFilterFields{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getAllSupplyOrderUsecase",
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					"/getAllSupplyOrderUsecase",
					nil,
				),
			)
		}
	}

	supplyOrders, err := getAllSOUcase.getAllSupplyOrderRepo.Execute(
		supplyorderdomainrepositorytypes.GetAllSupplyOrderInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllSupplyOrderUsecase",
			err,
		)
	}

	return supplyOrders, nil
}
