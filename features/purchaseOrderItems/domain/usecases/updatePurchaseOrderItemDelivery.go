package purchaseorderitempresentationusecases

import (
	"encoding/json"

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
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updatePurchaseOrderItemDeliveryUsecase struct {
	getAccountFromAuthDataRepo                    accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo                    memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatePurchaseOrderItemDeliveryRepo    purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemDeliveryRepository
	updatePurchaseOrderItemDeliveryAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                                  string
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
		"UpdatePurchaseOrderItemDeliveryUsecase",
	}, nil
}

func (updatePurchaseOrderItemDeliveryUcase *updatePurchaseOrderItemDeliveryUsecase) validation(input purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput) (purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput, error) {
	if &input.Context == nil {
		return purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updatePurchaseOrderItemDeliveryUcase.pathIdentity,
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
			updatePurchaseOrderItemDeliveryUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updatePurchaseOrderItemDeliveryUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := updatePurchaseOrderItemDeliveryUcase.getAccountMemberAccessRepo.Execute(
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
			updatePurchaseOrderItemDeliveryUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.PurchaseOrderItemDeliveryAccesses.PurchaseOrderItemDeliveryUpdate"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.PurchaseOrderItemDeliveryAccesses.PurchaseOrderItemDeliveryAssignDriver"), false,
		).(bool); accessible {
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updatePurchaseOrderItemDeliveryUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					updatePurchaseOrderItemDeliveryUcase.pathIdentity,
					nil,
				),
			)
		}
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
			updatePurchaseOrderItemDeliveryUcase.pathIdentity,
			err,
		)
	}

	return updatePurchaseOrderItemDeliveryOutput, nil
}
