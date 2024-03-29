package purchaseorderpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updatePurchaseOrderUsecase struct {
	getAccountFromAuthDataRepo        accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo        memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatePurchaseOrderRepo    purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository
	approveUpdatePurchaseOrderRepo    purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository
	updatePurchaseOrderByCronRepo     purchaseorderdomainrepositoryinterfaces.UpdatePurchaseOrderByCronRepository
	updatePurchaseOrderAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                      string
}

func NewUpdatePurchaseOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdatePurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository,
	approveUpdatePurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository,
	updatePurchaseOrderByCronRepo purchaseorderdomainrepositoryinterfaces.UpdatePurchaseOrderByCronRepository,
) (purchaseorderpresentationusecaseinterfaces.UpdatePurchaseOrderUsecase, error) {
	return &updatePurchaseOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdatePurchaseOrderRepo,
		approveUpdatePurchaseOrderRepo,
		updatePurchaseOrderByCronRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderAccesses: &model.PurchaseOrderAccessesInput{
				PurchaseOrderUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdatePurchaseOrderUsecase",
	}, nil
}

func (updatePurchaseOrderUcase *updatePurchaseOrderUsecase) validation(input purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput) (purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput, error) {
	if &input.Context == nil {
		return purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updatePurchaseOrderUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updatePurchaseOrderUcase *updatePurchaseOrderUsecase) Execute(input purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput) (*model.PurchaseOrder, error) {
	validatedInput, err := updatePurchaseOrderUcase.validation(input)
	if err != nil {
		return nil, err
	}

	if validatedInput.CronAuthenticated {
		_, err := updatePurchaseOrderUcase.updatePurchaseOrderByCronRepo.RunTransaction()
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updatePurchaseOrderUcase.pathIdentity,
				err,
			)
		}
		return nil, nil
	}

	purchaseOrderToUpdate := &model.InternalUpdatePurchaseOrder{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdatePurchaseOrder)
	json.Unmarshal(jsonTemp, purchaseOrderToUpdate)

	account, err := updatePurchaseOrderUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updatePurchaseOrderUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updatePurchaseOrderUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updatePurchaseOrderUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updatePurchaseOrderUcase.updatePurchaseOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updatePurchaseOrderUcase.pathIdentity,
			err,
		)
	}

	jsonTemp, _ = json.Marshal(accMemberAccess)
	json.Unmarshal(jsonTemp, &purchaseOrderToUpdate.MemberAccess)

	for i, poItem := range validatedInput.UpdatePurchaseOrder.Items {
		if poItem.PurchaseOrderItemReturn != nil {
			for j, descriptivePhoto := range poItem.PurchaseOrderItemReturn.Photos {
				purchaseOrderToUpdate.Items[i].PurchaseOrderItemReturn.Photos[j].Photo.File = descriptivePhoto.Photo.File
			}
		}
	}

	// if user is only going to approve proposal
	if purchaseOrderToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updatePurchaseOrderUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updatePurchaseOrderUcase.pathIdentity,
				nil,
			)
		}

		purchaseOrderToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updatePurchaseOrderOutput, err := updatePurchaseOrderUcase.approveUpdatePurchaseOrderRepo.RunTransaction(
			purchaseOrderToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updatePurchaseOrderUcase.pathIdentity,
				err,
			)
		}

		return updatePurchaseOrderOutput, nil
	}

	purchaseOrderToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval != nil {
		if *accMemberAccess.Access.PurchaseOrderAccesses.PurchaseOrderApproval {
			purchaseOrderToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	purchaseOrderToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updatePurchaseOrderOutput, err := updatePurchaseOrderUcase.proposeUpdatePurchaseOrderRepo.RunTransaction(
		purchaseOrderToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updatePurchaseOrderUcase.pathIdentity,
			err,
		)
	}

	return updatePurchaseOrderOutput, nil
}
