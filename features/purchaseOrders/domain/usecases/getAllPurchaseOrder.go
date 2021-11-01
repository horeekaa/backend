package purchaseorderpresentationusecases

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
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/types"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllPurchaseOrderUsecase struct {
	getAccountFromAuthDataRepo        accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo        memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllpurchaseOrderRepo           purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository
	getAllpurchaseOrderAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllPurchaseOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllpurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository,
) (purchaseorderpresentationusecaseinterfaces.GetAllPurchaseOrderUsecase, error) {
	return &getAllPurchaseOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllpurchaseOrderRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllPOUcase *getAllPurchaseOrderUsecase) validation(input purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput) (*purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput, error) {
	if &input.Context == nil {
		return &purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllPurchaseOrderUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllPOUcase *getAllPurchaseOrderUsecase) Execute(
	input purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput,
) ([]*model.PurchaseOrder, error) {
	validatedInput, err := getAllPOUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllPOUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllPurchaseOrderUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllPOUcase.getAccountMemberAccessRepo.Execute(
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
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.PurchaseOrderAccesses.PurchaseOrderReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.PurchaseOrderAccesses.PurchaseOrderReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.Organization = &model.OrganizationFilterFields{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getAllPurchaseOrderUsecase",
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					"/getAllPurchaseOrderUsecase",
					nil,
				),
			)
		}
	}

	purchaseOrders, err := getAllPOUcase.getAllpurchaseOrderRepo.Execute(
		purchaseorderdomainrepositorytypes.GetAllPurchaseOrderInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderUsecase",
			err,
		)
	}

	return purchaseOrders, nil
}
