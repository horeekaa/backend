package paymentdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createPaymentRepository struct {
	createPaymentTransactionComponent paymentdomainrepositoryinterfaces.CreatePaymentTransactionComponent
	createDescriptivePhotoComponent   descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	updateInvoiceTrxComponent         invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreatePaymentRepository(
	createPaymentRepositoryTransactionComponent paymentdomainrepositoryinterfaces.CreatePaymentTransactionComponent,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (paymentdomainrepositoryinterfaces.CreatePaymentRepository, error) {
	createPaymentRepo := &createPaymentRepository{
		createPaymentRepositoryTransactionComponent,
		createDescriptivePhotoComponent,
		updateInvoiceTrxComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createPaymentRepo,
		"CreatePaymentRepository",
	)

	return createPaymentRepo, nil
}

func (createPaymentRepo *createPaymentRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createPaymentRepo.createPaymentTransactionComponent.PreTransaction(
		input.(*model.InternalCreatePayment),
	)
}

func (createPaymentRepo *createPaymentRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	paymentToCreate := input.(*model.InternalCreatePayment)
	generatedObjectID := createPaymentRepo.createPaymentTransactionComponent.GenerateNewObjectID()
	if paymentToCreate.Photo != nil {
		paymentToCreate.Photo.Category = model.DescriptivePhotoCategoryPurchaseOrderPaymentProof
		paymentToCreate.Photo.Object = &model.ObjectIDOnly{
			ID: &generatedObjectID,
		}
		paymentToCreate.Photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*paymentToCreate.ProposalStatus)
		paymentToCreate.Photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*paymentToCreate.SubmittingAccount)
		createdPhotoOutput, err := createPaymentRepo.createDescriptivePhotoComponent.TransactionBody(
			operationOption,
			paymentToCreate.Photo,
		)
		if err != nil {
			return nil, err
		}
		jsonTemp, _ := json.Marshal(createdPhotoOutput)
		json.Unmarshal(jsonTemp, &paymentToCreate.Photo)
	}

	if *paymentToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		_, err := createPaymentRepo.updateInvoiceTrxComponent.TransactionBody(
			operationOption,
			&model.InternalUpdateInvoice{
				ID: *paymentToCreate.Invoice.ID,
				Payments: []*model.ObjectIDOnly{
					{ID: &generatedObjectID},
				},
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return createPaymentRepo.createPaymentTransactionComponent.TransactionBody(
		operationOption,
		paymentToCreate,
	)
}

func (createPaymentRepo *createPaymentRepository) RunTransaction(
	input *model.InternalCreatePayment,
) (*model.Payment, error) {
	output, err := createPaymentRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Payment), err
}
