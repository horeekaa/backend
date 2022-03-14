package paymentdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePaymentRepository struct {
	paymentDataSource                        databasepaymentdatasourceinterfaces.PaymentDataSource
	approveUpdateDescriptivePhotoComponent   descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	approveUpdatePaymentTransactionComponent paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent
	updateInvoiceTrxComponent                invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                             string
}

func NewApproveUpdatePaymentRepository(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	approveUpdatePaymentTransactionComponent paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (paymentdomainrepositoryinterfaces.ApproveUpdatePaymentRepository, error) {
	approveUpdatePaymentRepo := &approveUpdatePaymentRepository{
		paymentDataSource,
		approveUpdateDescriptivePhotoComponent,
		approveUpdatePaymentTransactionComponent,
		updateInvoiceTrxComponent,
		mongoDBTransaction,
		"ApproveUpdatePaymentRepository",
	}

	mongoDBTransaction.SetTransaction(
		approveUpdatePaymentRepo,
		"ApproveUpdatePaymentRepository",
	)

	return approveUpdatePaymentRepo, nil
}

func (approveUpdatePaymentRepo *approveUpdatePaymentRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdatePaymentRepo *approveUpdatePaymentRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	paymentToApprove := input.(*model.InternalUpdatePayment)
	existingPayment, err := approveUpdatePaymentRepo.paymentDataSource.GetMongoDataSource().FindByID(
		paymentToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdatePaymentRepo.pathIdentity,
			err,
		)
	}
	if existingPayment.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
			ID: &existingPayment.Photo.ID,
		}
		updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*paymentToApprove.RecentApprovingAccount)
		updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*paymentToApprove.ProposalStatus)

		_, err := approveUpdatePaymentRepo.approveUpdateDescriptivePhotoComponent.TransactionBody(
			operationOption,
			updateDescriptivePhoto,
		)
		if err != nil {
			return nil, err
		}
	}

	if paymentToApprove.ProposalStatus != nil {
		if *paymentToApprove.ProposalStatus == model.EntityProposalStatusApproved {
			_, err := approveUpdatePaymentRepo.updateInvoiceTrxComponent.TransactionBody(
				operationOption,
				&model.InternalUpdateInvoice{
					ID: existingPayment.Invoice.ID,
					Payments: []*model.ObjectIDOnly{
						{ID: &existingPayment.ID},
					},
				},
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return approveUpdatePaymentRepo.approveUpdatePaymentTransactionComponent.TransactionBody(
		operationOption,
		paymentToApprove,
	)
}

func (approveUpdatePaymentRepo *approveUpdatePaymentRepository) RunTransaction(
	input *model.InternalUpdatePayment,
) (*model.Payment, error) {
	output, err := approveUpdatePaymentRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Payment), err
}
