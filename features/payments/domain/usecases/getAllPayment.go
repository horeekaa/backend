package paymentpresentationusecases

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
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentdomainrepositorytypes "github.com/horeekaa/backend/features/payments/domain/repositories/types"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllPaymentUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllPaymentRepo           paymentdomainrepositoryinterfaces.GetAllPaymentRepository
	getAllPaymentAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewGetAllPaymentUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllPaymentRepo paymentdomainrepositoryinterfaces.GetAllPaymentRepository,
) (paymentpresentationusecaseinterfaces.GetAllPaymentUsecase, error) {
	return &getAllPaymentUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllPaymentRepo,
		&model.MemberAccessRefOptionsInput{
			PaymentAccesses: &model.PaymentAccessesInput{
				PaymentReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllPaymentUsecase",
	}, nil
}

func (getAllPaymentUcase *getAllPaymentUsecase) validation(input paymentpresentationusecasetypes.GetAllPaymentUsecaseInput) (*paymentpresentationusecasetypes.GetAllPaymentUsecaseInput, error) {
	if &input.Context == nil {
		return &paymentpresentationusecasetypes.GetAllPaymentUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllPaymentUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllPaymentUcase *getAllPaymentUsecase) Execute(
	input paymentpresentationusecasetypes.GetAllPaymentUsecaseInput,
) ([]*model.Payment, error) {
	validatedInput, err := getAllPaymentUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllPaymentUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllPaymentUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllPaymentUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllPaymentUcase.getAccountMemberAccessRepo.Execute(
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
			getAllPaymentUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.PaymentAccesses.PaymentReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.PaymentAccesses.PaymentReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.Organization = &model.OrganizationForPaymentFilterFields{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAllPaymentUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					getAllPaymentUcase.pathIdentity,
					nil,
				),
			)
		}
	}

	payments, err := getAllPaymentUcase.getAllPaymentRepo.Execute(
		paymentdomainrepositorytypes.GetAllPaymentInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllPaymentUcase.pathIdentity,
			err,
		)
	}

	return payments, nil
}
