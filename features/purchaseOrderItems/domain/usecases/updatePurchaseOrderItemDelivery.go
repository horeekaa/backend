package purchaseorderitempresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updatePurchaseOrderItemDeliveryUsecase struct {
	getAccountFromAuthDataRepo                    accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo                    memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatePurchaseOrderItemDeliveryRepo    purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemDeliveryRepository
	updatePurchaseOrderItemDeliveryAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdatePurchaseOrderItemDeliveryUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdatePurchaseOrderItemDeliveryRepo purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemDeliveryRepository,
) (purchaseorderitempresentationusecaseinterfaces.UpdatePurchaseOrderItemDeliveryUsecase, error) {
	return &updatePurchaseOrderItemDeliveryUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdatePurchaseOrderItemDeliveryRepo,
		&model.MemberAccessRefOptionsInput{
			PurchaseOrderItemDeliveryAccesses: &model.PurchaseOrderItemDeliveryAccessesInput{
				PurchaseOrderItemDeliveryUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updatePurchaseOrderItemDeliveryUcase *updatePurchaseOrderItemDeliveryUsecase) validation(input purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput) (purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput, error) {
	if &input.Context == nil {
		return purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updatePurchaseOrderItemDeliveryUsecase",
				nil,
			)
	}

	return input, nil
}

func (updatePurchaseOrderItemDeliveryUcase *updatePurchaseOrderItemDeliveryUsecase) Execute(input purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput) (*model.PurchaseOrderItem, error) {
	validatedInput, err := updatePurchaseOrderItemDeliveryUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updatePurchaseOrderItemDeliveryUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePurchaseOrderItemDeliveryUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updatePurchaseOrderItemDeliveryUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = updatePurchaseOrderItemDeliveryUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updatePurchaseOrderItemDeliveryUcase.updatePurchaseOrderItemDeliveryAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePurchaseOrderItemDeliveryUsecase",
			err,
		)
	}

	purchaseOrderItemToUpdate := &model.InternalUpdatePurchaseOrderItem{
		ID: &validatedInput.UpdatePurchaseOrderItemDelivery.ID,
	}
	jsonTemp, _ := json.Marshal(validatedInput.UpdatePurchaseOrderItemDelivery)
	json.Unmarshal(jsonTemp, purchaseOrderItemToUpdate)

	for i, descriptivePhoto := range validatedInput.UpdatePurchaseOrderItemDelivery.DeliveryDetail.Photos {
		if descriptivePhoto.Photo != nil {
			purchaseOrderItemToUpdate.DeliveryDetail.Photos[i].Photo.File = descriptivePhoto.Photo.File
		}
	}
	for i, descriptivePhoto := range validatedInput.UpdatePurchaseOrderItemDelivery.DeliveryDetail.PhotosAfterReceived {
		if descriptivePhoto.Photo != nil {
			purchaseOrderItemToUpdate.DeliveryDetail.PhotosAfterReceived[i].Photo.File = descriptivePhoto.Photo.File
		}
	}

	purchaseOrderItemToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)

	purchaseOrderItemToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updatePurchaseOrderItemDeliveryOutput, err := updatePurchaseOrderItemDeliveryUcase.proposeUpdatePurchaseOrderItemDeliveryRepo.RunTransaction(
		purchaseOrderItemToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updatePurchaseOrderItemDeliveryUsecase",
			err,
		)
	}

	return updatePurchaseOrderItemDeliveryOutput, nil
}
