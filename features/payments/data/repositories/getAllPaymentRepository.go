package paymentdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentdomainrepositorytypes "github.com/horeekaa/backend/features/payments/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllPaymentRepository struct {
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity      string
}

func NewGetAllPaymentRepository(
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (paymentdomainrepositoryinterfaces.GetAllPaymentRepository, error) {
	return &getAllPaymentRepository{
		paymentDataSource,
		mongoQueryBuilder,
		"GetAllPaymentRepository",
	}, nil
}

func (getAllPaymentRepo *getAllPaymentRepository) Execute(
	input paymentdomainrepositorytypes.GetAllPaymentInput,
) ([]*model.Payment, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllPaymentRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	payments, err := getAllPaymentRepo.paymentDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllPaymentRepo.pathIdentity,
			err,
		)
	}

	return payments, nil
}
