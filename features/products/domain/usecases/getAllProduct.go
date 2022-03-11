package productpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productdomainrepositorytypes "github.com/horeekaa/backend/features/products/domain/repositories/types"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllProductUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllProductRepo           productdomainrepositoryinterfaces.GetAllProductRepository
	getAllProductAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewGetAllProductUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllProductRepo productdomainrepositoryinterfaces.GetAllProductRepository,
) (productpresentationusecaseinterfaces.GetAllProductUsecase, error) {
	return &getAllProductUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllProductRepo,
		&model.MemberAccessRefOptionsInput{
			ProductAccesses: &model.ProductAccessesInput{
				ProductReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllProductUsecase",
	}, nil
}

func (getAllProdUcase *getAllProductUsecase) validation(input productpresentationusecasetypes.GetAllProductUsecaseInput) (*productpresentationusecasetypes.GetAllProductUsecaseInput, error) {
	if &input.Context == nil {
		return &productpresentationusecasetypes.GetAllProductUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllProdUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllProdUcase *getAllProductUsecase) Execute(
	input productpresentationusecasetypes.GetAllProductUsecaseInput,
) ([]*model.Product, error) {
	validatedInput, err := getAllProdUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllProdUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllProdUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllProdUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllProdUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllProdUcase.getAllProductAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllProdUcase.pathIdentity,
			err,
		)
	}

	products, err := getAllProdUcase.getAllProductRepo.Execute(
		productdomainrepositorytypes.GetAllProductInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllProdUcase.pathIdentity,
			err,
		)
	}

	return products, nil
}
