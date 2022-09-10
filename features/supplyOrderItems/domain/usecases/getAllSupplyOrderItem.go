package supplyorderitempresentationusecases

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
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/types"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllSupplyOrderItemUsecase struct {
	getAccountFromAuthDataRepo    accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo    memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	supplyOrderItemRepo           supplyorderitemdomainrepositoryinterfaces.GetAllSupplyOrderItemRepository
	supplyOrderItemAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                  string
}

func NewGetAllSupplyOrderItemUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	supplyOrderItemRepo supplyorderitemdomainrepositoryinterfaces.GetAllSupplyOrderItemRepository,
) (supplyorderitempresentationusecaseinterfaces.GetAllSupplyOrderItemUsecase, error) {
	return &getAllSupplyOrderItemUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		supplyOrderItemRepo,
		&model.MemberAccessRefOptionsInput{
			SupplyOrderAccesses: &model.SupplyOrderAccessesInput{
				SupplyOrderReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllSupplyOrderItemUsecase",
	}, nil
}

func (getAllSOItemUcase *getAllSupplyOrderItemUsecase) validation(input supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput) (*supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput, error) {
	if &input.Context == nil {
		return &supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllSOItemUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllSOItemUcase *getAllSupplyOrderItemUsecase) Execute(
	input supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput,
) ([]*model.SupplyOrderItem, error) {
	validatedInput, err := getAllSOItemUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllSOItemUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllSOItemUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllSOItemUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllSOItemUcase.getAccountMemberAccessRepo.Execute(
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
			getAllSOItemUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.SupplyOrderAccesses.SupplyOrderReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.SupplyOrderAccesses.SupplyOrderReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.SupplyOrder = &model.SupplyOrderForSupplyOrderItemFilterFields{
				Organization: &model.ObjectIDOnlyFilterField{
					ID: &model.ObjectIDFilterField{
						Value:     &memberAccess.Organization.ID,
						Operation: model.ObjectIDOperationEqual,
					},
				},
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAllSOItemUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					getAllSOItemUcase.pathIdentity,
					nil,
				),
			)
		}
	}

	supplyOrderItems, err := getAllSOItemUcase.supplyOrderItemRepo.Execute(
		supplyorderitemdomainrepositorytypes.GetAllSupplyOrderItemInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllSOItemUcase.pathIdentity,
			err,
		)
	}

	return supplyOrderItems, nil
}
