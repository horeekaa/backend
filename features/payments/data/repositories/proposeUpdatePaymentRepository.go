package paymentdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePaymentRepository struct {
	paymentDataSource                        databasepaymentdatasourceinterfaces.PaymentDataSource
	createDescriptivePhotoComponent          descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	proposeUpdateDescriptivePhotoComponent   descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	proposeUpdatePaymentTransactionComponent paymentdomainrepositoryinterfaces.ProposeUpdatePaymentTransactionComponent
	updateInvoiceTrxComponent                invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                             string
}

func NewProposeUpdatePaymentRepository(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	proposeUpdatePaymentRepositoryTransactionComponent paymentdomainrepositoryinterfaces.ProposeUpdatePaymentTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository, error) {
	proposeUpdatePaymentRepo := &proposeUpdatePaymentRepository{
		paymentDataSource,
		createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent,
		proposeUpdatePaymentRepositoryTransactionComponent,
		updateInvoiceTrxComponent,
		mongoDBTransaction,
		"ProposeUpdatePaymentRepository",
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdatePaymentRepo,
		"ProposeUpdatePaymentRepository",
	)

	return proposeUpdatePaymentRepo, nil
}

func (updatePaymentRepo *proposeUpdatePaymentRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updatePaymentRepo.proposeUpdatePaymentTransactionComponent.PreTransaction(
		input.(*model.InternalUpdatePayment),
	)
}

func (updatePaymentRepo *proposeUpdatePaymentRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	paymentToUpdate := input.(*model.InternalUpdatePayment)
	existingPayment, err := updatePaymentRepo.paymentDataSource.GetMongoDataSource().FindByID(
		*paymentToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePaymentRepo.pathIdentity,
			err,
		)
	}

	if paymentToUpdate.Photo != nil {
		if paymentToUpdate.Photo.ID != nil {
			paymentToUpdate.Photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*paymentToUpdate.ProposalStatus)
			paymentToUpdate.Photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*paymentToUpdate.SubmittingAccount)

			_, err := updatePaymentRepo.proposeUpdateDescriptivePhotoComponent.TransactionBody(
				operationOption,
				paymentToUpdate.Photo,
			)
			if err != nil {
				return nil, err
			}
		} else {
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(paymentToUpdate.Photo)
			json.Unmarshal(jsonTemp, photoToCreate)
			if funk.Get(paymentToUpdate, "Photo.Photo") != nil {
				photoToCreate.Photo.File = paymentToUpdate.Photo.Photo.File
			}
			photoToCreate.Category = model.DescriptivePhotoCategoryPaymentProof
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingPayment.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*paymentToUpdate.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*paymentToUpdate.SubmittingAccount)

			createdDescriptivePhoto, err := updatePaymentRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photoToCreate,
			)
			if err != nil {
				return nil, err
			}

			if existingPayment.Photo != nil {
				_, err = updatePaymentRepo.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					operationOption,
					&model.InternalUpdateDescriptivePhoto{
						ID:       &existingPayment.Photo.ID,
						IsActive: func(b bool) *bool { return &b }(false),
						SubmittingAccount: func(m model.ObjectIDOnly) *model.ObjectIDOnly {
							return &m
						}(*paymentToUpdate.SubmittingAccount),
						ProposalStatus: func(s model.EntityProposalStatus) *model.EntityProposalStatus {
							return &s
						}(*paymentToUpdate.ProposalStatus),
					},
				)
				if err != nil {
					return nil, err
				}
			}
			paymentToUpdate.Photo = &model.InternalUpdateDescriptivePhoto{
				ID: &createdDescriptivePhoto.ID,
			}
		}
	}

	updatedPayment, err := updatePaymentRepo.proposeUpdatePaymentTransactionComponent.TransactionBody(
		operationOption,
		paymentToUpdate,
	)
	if err != nil {
		return nil, err
	}

	if paymentToUpdate.ProposalStatus != nil {
		if updatedPayment.ProposedChanges.Invoice != nil {
			_, err := updatePaymentRepo.updateInvoiceTrxComponent.TransactionBody(
				operationOption,
				&model.InternalUpdateInvoice{
					ID: updatedPayment.ProposedChanges.Invoice.ID,
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

	return updatedPayment, nil
}

func (updatePaymentRepo *proposeUpdatePaymentRepository) RunTransaction(
	input *model.InternalUpdatePayment,
) (*model.Payment, error) {
	output, err := updatePaymentRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Payment), err
}
