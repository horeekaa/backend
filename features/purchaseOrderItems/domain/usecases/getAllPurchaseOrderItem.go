package purchaseorderitempresentationusecases

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
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/types"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllPurchaseOrderItemUsecase struct {
	getAccountFromAuthDataRepo            accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo            memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllPurchaseOrderItemRepo           purchaseorderitemdomainrepositoryinterfaces.GetAllPurchaseOrderItemRepository
	getAllPurchaseOrderItemAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllPurchaseOrderItemUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllPurchaseOrderItemRepo purchaseorderitemdomainrepositoryinterfaces.GetAllPurchaseOrderItemRepository,
) (purchaseorderitempresentationusecaseinterfaces.GetAllPurchaseOrderItemUsecase, error) {
	return &getAllPurchaseOrderItemUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllPurchaseOrderItemRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllPOItemUcase *getAllPurchaseOrderItemUsecase) validation(input purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput) (*purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput, error) {
	if &input.Context == nil {
		return &purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllPurchaseOrderItemUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllPOItemUcase *getAllPurchaseOrderItemUsecase) Execute(
	input purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput,
) ([]*model.PurchaseOrderItem, error) {
	validatedInput, err := getAllPOItemUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllPOItemUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderItemUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllPurchaseOrderItemUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllPOItemUcase.getAccountMemberAccessRepo.Execute(
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
			"/getAllPurchaseOrderItemUsecase",
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.PurchaseOrderAccesses.PurchaseOrderReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.PurchaseOrderAccesses.PurchaseOrderItemReadOwned"), false,
		).(bool); accessible {
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getAllPurchaseOrderItemUsecase",
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					"/getAllPurchaseOrderItemUsecase",
					nil,
				),
			)
		}
	}

	purchaseOrderItems, err := getAllPOItemUcase.getAllPurchaseOrderItemRepo.Execute(
		purchaseorderitemdomainrepositorytypes.GetAllPurchaseOrderItemInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllPurchaseOrderItemUsecase",
			err,
		)
	}

	return purchaseOrderItems, nil
}
