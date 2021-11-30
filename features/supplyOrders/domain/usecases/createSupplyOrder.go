package supplyorderpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createSupplyOrderUsecase struct {
	getAccountFromAuthDataRepo      accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo      memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createSupplyOrderRepo           supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository
	createSupplyOrderAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateSupplyOrderUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createSupplyOrderRepo supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository,
) (supplyorderpresentationusecaseinterfaces.CreateSupplyOrderUsecase, error) {
	return &createSupplyOrderUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createSupplyOrderRepo,
		&model.MemberAccessRefOptionsInput{
			SupplyOrderAccesses: &model.SupplyOrderAccessesInput{
				SupplyOrderCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createSupplyOrderUcase *createSupplyOrderUsecase) validation(input supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput) (supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput, error) {
	if &input.Context == nil {
		return supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createSupplyOrderUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateSupplyOrder.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createSupplyOrderUcase *createSupplyOrderUsecase) Execute(input supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput) ([]*model.SupplyOrder, error) {
	validatedInput, err := createSupplyOrderUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createSupplyOrderUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createSupplyOrderUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createSupplyOrderUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createSupplyOrderUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createSupplyOrderUcase.createSupplyOrderAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createSupplyOrderUsecase",
			err,
		)
	}
	if accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval != nil {
		if *accMemberAccess.Access.SupplyOrderAccesses.SupplyOrderApproval {
			validatedInput.CreateSupplyOrder.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	supplyOrderToCreate := &model.InternalCreateSupplyOrder{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateSupplyOrder)
	json.Unmarshal(jsonTemp, supplyOrderToCreate)

	jsonTemp, _ = json.Marshal(accMemberAccess)
	json.Unmarshal(jsonTemp, &supplyOrderToCreate.MemberAccess)

	for i, soItem := range validatedInput.CreateSupplyOrder.Items {
		if soItem.PickUpDetail != nil {
			for j, descriptivePhoto := range soItem.PickUpDetail.Photos {
				supplyOrderToCreate.Items[i].PickUpDetail.Photos[j].Photo.File = descriptivePhoto.Photo.File
			}
		}
	}

	supplyOrderToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdSupplyOrder, err := createSupplyOrderUcase.createSupplyOrderRepo.RunTransaction(
		supplyOrderToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createSupplyOrderUsecase",
			err,
		)
	}

	return createdSupplyOrder, nil
}
