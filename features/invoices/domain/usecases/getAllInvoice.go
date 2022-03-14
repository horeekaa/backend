package invoicepresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type getAllInvoiceUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllInvoiceRepo           invoicedomainrepositoryinterfaces.GetAllInvoiceRepository
	getAllInvoiceAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewGetAllInvoiceUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllInvoiceRepo invoicedomainrepositoryinterfaces.GetAllInvoiceRepository,
) (invoicepresentationusecaseinterfaces.GetAllInvoiceUsecase, error) {
	return &getAllInvoiceUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllInvoiceRepo,
		&model.MemberAccessRefOptionsInput{
			InvoiceAccesses: &model.InvoiceAccessesInput{
				InvoiceReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllInvoiceUsecase",
	}, nil
}

func (getAllInvoiceUcase *getAllInvoiceUsecase) validation(input invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput) (*invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput, error) {
	if &input.Context == nil {
		return &invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllInvoiceUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllInvoiceUcase *getAllInvoiceUsecase) Execute(
	input invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput,
) ([]*model.Invoice, error) {
	validatedInput, err := getAllInvoiceUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllInvoiceUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllInvoiceUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllInvoiceUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	memberAccess, err := getAllInvoiceUcase.getAccountMemberAccessRepo.Execute(
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
			getAllInvoiceUcase.pathIdentity,
			err,
		)
	}
	if accessible := funk.GetOrElse(
		funk.Get(memberAccess, "Access.InvoiceAccesses.InvoiceReadAll"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(memberAccess, "Access.InvoiceAccesses.InvoiceReadOwned"), false,
		).(bool); accessible {
			validatedInput.FilterFields.Organization = &model.OrganizationForInvoiceFilterFields{
				ID: &memberAccess.Organization.ID,
			}
		} else {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAllInvoiceUcase.pathIdentity,
				horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					getAllInvoiceUcase.pathIdentity,
					nil,
				),
			)
		}
	}

	invoices, err := getAllInvoiceUcase.getAllInvoiceRepo.Execute(
		invoicedomainrepositorytypes.GetAllInvoiceInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllInvoiceUcase.pathIdentity,
			err,
		)
	}

	return invoices, nil
}
