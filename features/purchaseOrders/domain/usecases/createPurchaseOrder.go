package purchaseorderpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderUsecase struct {
	getAccountFromAuthDataRepo        accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo        memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createPurchaseOrderRepo           purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository
	createPurchaseOrderAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreatePurchaseOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createpurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository,
) (purchaseorderpresentationusecaseinterfaces.CreatePurchaseOrderUsecase, error) {
	return &createPurchaseOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createpurchaseOrderRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createPurchaseOrderUcase *createPurchaseOrderUsecase) validation(input purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput) (purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput, error) {
	if &input.Context == nil {
		return purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createPurchaseOrderUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreatePurchaseOrder.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createPurchaseOrderUcase *createPurchaseOrderUsecase) Execute(input purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput) ([]*model.PurchaseOrder, error) {
	validatedInput, err := createPurchaseOrderUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createPurchaseOrderUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPurchaseOrderUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createPurchaseOrderUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createPurchaseOrderUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createPurchaseOrderUcase.createPurchaseOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPurchaseOrderUsecase",
			err,
		)
	}
	if accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval != nil {
		if *accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval {
			validatedInput.CreatePurchaseOrder.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	purchaseOrderToCreate := &model.InternalCreatePurchaseOrder{}
	jsonTemp, _ := json.Marshal(validatedInput.CreatePurchaseOrder)
	json.Unmarshal(jsonTemp, purchaseOrderToCreate)

	purchaseOrderToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdPurchaseOrder, err := createPurchaseOrderUcase.createPurchaseOrderRepo.RunTransaction(
		purchaseOrderToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createPurchaseOrderUsecase",
			err,
		)
	}

	return createdPurchaseOrder, nil
}