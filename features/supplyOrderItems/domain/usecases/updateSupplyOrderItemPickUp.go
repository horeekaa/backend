package supplyorderitempresentationusecases

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
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updateSupplyOrderItemPickUpUsecase struct {
	getAccountFromAuthDataRepo                accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo                memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdatesupplyOrderItemPickUpRepo    supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository
	updatesupplyOrderItemPickUpAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                              string
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
		"UpdateSupplyOrderItemPickUpUsecase",
	}, nil
}

func (updateSupplyOrderItemPickUpUcase *updateSupplyOrderItemPickUpUsecase) validation(input supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput) (supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput, error) {
	if &input.Context == nil {
		return supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateSupplyOrderItemPickUpUcase.pathIdentity,
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
			updateSupplyOrderItemPickUpUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateSupplyOrderItemPickUpUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := updateSupplyOrderItemPickUpUcase.getAccountMemberAccessRepo.Execute(
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
			updateSupplyOrderItemPickUpUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.SupplyOrderItemPickUpAccesses.SupplyOrderItemPickUpUpdate"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.SupplyOrderItemPickUpAccesses.SupplyOrderItemPickUpAssignDriver"), false,
		).(bool); accessible {
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateSupplyOrderItemPickUpUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					updateSupplyOrderItemPickUpUcase.pathIdentity,
					nil,
				),
			)
		}
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
			updateSupplyOrderItemPickUpUcase.pathIdentity,
			err,
		)
	}

	return updateSupplyOrderItemPickUpOutput, nil
}
