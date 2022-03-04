package paymentpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getPaymentUsecase struct {
	getPaymentRepository paymentdomainrepositoryinterfaces.GetPaymentRepository
}

func NewGetPaymentUsecase(
	getPaymentRepository paymentdomainrepositoryinterfaces.GetPaymentRepository,
) (paymentpresentationusecaseinterfaces.GetPaymentUsecase, error) {
	return &getPaymentUsecase{
		getPaymentRepository,
	}, nil
}

func (getPaymentUcase *getPaymentUsecase) validation(
	input *model.PaymentFilterFields,
) (*model.PaymentFilterFields, error) {
	return input, nil
}

func (getPaymentUcase *getPaymentUsecase) Execute(
	filterFields *model.PaymentFilterFields,
) (*model.Payment, error) {
	validatedFilterFields, err := getPaymentUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	payment, err := getPaymentUcase.getPaymentRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPayment",
			err,
		)
	}
	return payment, nil
}
