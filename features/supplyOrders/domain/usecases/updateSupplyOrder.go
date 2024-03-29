package supplyorderpresentationusecases

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
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updateSupplyOrderUsecase struct {
	getAccountFromAuthDataRepo      accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo      memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatesupplyOrderRepo    supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderRepository
	approveUpdatesupplyOrderRepo    supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderRepository
	updatesupplyOrderAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                    string
}

func NewUpdateSupplyOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdatesupplyOrderRepo supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderRepository,
	approveUpdatesupplyOrderRepo supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderRepository,
) (supplyorderpresentationusecaseinterfaces.UpdateSupplyOrderUsecase, error) {
	return &updateSupplyOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdatesupplyOrderRepo,
		approveUpdatesupplyOrderRepo,
		&model.MemberAccessRefOptionsInput{
			SupplyOrderAccesses: &model.SupplyOrderAccessesInput{
				SupplyOrderUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateSupplyOrderUsecase",
	}, nil
}

func (updateSupplyOrderUcase *updateSupplyOrderUsecase) validation(input supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput) (supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput, error) {
	if &input.Context == nil {
		return supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateSupplyOrderUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateSupplyOrderUcase *updateSupplyOrderUsecase) Execute(input supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput) (*model.SupplyOrder, error) {
	validatedInput, err := updateSupplyOrderUcase.validation(input)
	if err != nil {
		return nil, err
	}
	supplyOrderToUpdate := &model.InternalUpdateSupplyOrder{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateSupplyOrder)
	json.Unmarshal(jsonTemp, supplyOrderToUpdate)

	account, err := updateSupplyOrderUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateSupplyOrderUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateSupplyOrderUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateSupplyOrderUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateSupplyOrderUcase.updatesupplyOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateSupplyOrderUcase.pathIdentity,
			err,
		)
	}

	jsonTemp, _ = json.Marshal(accMemberAccess)
	json.Unmarshal(jsonTemp, &supplyOrderToUpdate.MemberAccess)

	for i, soItem := range validatedInput.UpdateSupplyOrder.Items {
		for j, descriptivePhoto := range soItem.Photos {
			supplyOrderToUpdate.Items[i].Photos[j].Photo.File = descriptivePhoto.Photo.File
		}

		if soItem.SupplyOrderItemReturn != nil {
			for j, descriptivePhoto := range soItem.SupplyOrderItemReturn.Photos {
				supplyOrderToUpdate.Items[i].SupplyOrderItemReturn.Photos[j].Photo.File = descriptivePhoto.Photo.File
			}
		}
	}

	if funk.Get(validatedInput.UpdateSupplyOrder, "Payment.Photo.Photo") != nil {
		supplyOrderToUpdate.Payment.Photo.Photo.File = validatedInput.UpdateSupplyOrder.Payment.Photo.Photo.File
	}

	// if user is only going to approve proposal
	if supplyOrderToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateSupplyOrderUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateSupplyOrderUcase.pathIdentity,
				nil,
			)
		}

		supplyOrderToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateSupplyOrderOutput, err := updateSupplyOrderUcase.approveUpdatesupplyOrderRepo.RunTransaction(
			supplyOrderToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateSupplyOrderUcase.pathIdentity,
				err,
			)
		}

		return updateSupplyOrderOutput, nil
	}

	supplyOrderToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval != nil {
		if *accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval {
			supplyOrderToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	supplyOrderToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateSupplyOrderOutput, err := updateSupplyOrderUcase.proposeUpdatesupplyOrderRepo.RunTransaction(
		supplyOrderToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateSupplyOrderUcase.pathIdentity,
			err,
		)
	}

	return updateSupplyOrderOutput, nil
}
