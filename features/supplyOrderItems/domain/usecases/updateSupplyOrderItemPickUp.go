package supplyorderitempresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type updateSupplyOrderItemPickUpUsecase struct {
	getAccountFromAuthDataRepo                accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo                memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatesupplyOrderItemPickUpRepo    supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository
	updatesupplyOrderItemPickUpAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewUpdateSupplyOrderItemPickUpUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdatesupplyOrderItemPickUpRepo supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository,
) (supplyorderitempresentationusecaseinterfaces.UpdateSupplyOrderItemPickUpUsecase, error) {
	return &updateSupplyOrderItemPickUpUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdatesupplyOrderItemPickUpRepo,
		&model.MemberAccessRefOptionsInput{
			SupplyOrderItemPickUpAccesses: &model.SupplyOrderItemPickUpAccessesInput{
				SupplyOrderItemPickUpUpdate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (updateSupplyOrderItemPickUpUcase *updateSupplyOrderItemPickUpUsecase) validation(input supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput) (supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput, error) {
	if &input.Context == nil {
		return supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/updateSupplyOrderItemPickUpUsecase",
				nil,
			)
	}

	return input, nil
}

func (updateSupplyOrderItemPickUpUcase *updateSupplyOrderItemPickUpUsecase) Execute(input supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput) (*model.SupplyOrderItem, error) {
	validatedInput, err := updateSupplyOrderItemPickUpUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := updateSupplyOrderItemPickUpUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateSupplyOrderItemPickUpUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/updateSupplyOrderItemPickUpUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = updateSupplyOrderItemPickUpUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateSupplyOrderItemPickUpUcase.updatesupplyOrderItemPickUpAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateSupplyOrderItemPickUpUsecase",
			err,
		)
	}

	supplyOrderItemToUpdate := &model.InternalUpdateSupplyOrderItem{
		ID: &validatedInput.UpdateSupplyOrderItemPickUp.ID,
	}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateSupplyOrderItemPickUp)
	json.Unmarshal(jsonTemp, supplyOrderItemToUpdate)

	for i, descriptivePhoto := range validatedInput.UpdateSupplyOrderItemPickUp.PickUpDetail.Photos {
		if descriptivePhoto.Photo != nil {
			supplyOrderItemToUpdate.PickUpDetail.Photos[i].Photo.File = descriptivePhoto.Photo.File
		}
	}

	supplyOrderItemToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)

	supplyOrderItemToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateSupplyOrderItemPickUpOutput, err := updateSupplyOrderItemPickUpUcase.proposeUpdatesupplyOrderItemPickUpRepo.RunTransaction(
		supplyOrderItemToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/updateSupplyOrderItemPickUpUsecase",
			err,
		)
	}

	return updateSupplyOrderItemPickUpOutput, nil
}
