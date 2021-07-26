package productpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type createProductUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	createProductRepo           productdomainrepositoryinterfaces.CreateProductRepository
	createproductAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewCreateProductUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	createProductRepo productdomainrepositoryinterfaces.CreateProductRepository,
) (productpresentationusecaseinterfaces.CreateProductUsecase, error) {
	return &createProductUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		createProductRepo,
		&model.MemberAccessRefOptionsInput{
			ProductAccesses: &model.ProductAccessesInput{
				ProductCreate: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (createProductUcase *createProductUsecase) validation(input productpresentationusecasetypes.CreateProductUsecaseInput) (productpresentationusecasetypes.CreateProductUsecaseInput, error) {
	if &input.Context == nil {
		return productpresentationusecasetypes.CreateProductUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/createProductUsecase",
				nil,
			)
	}
	proposedProposalStatus := model.EntityProposalStatusProposed
	input.CreateProduct.ProposalStatus = &proposedProposalStatus
	return input, nil
}

func (createProductUcase *createProductUsecase) Execute(input productpresentationusecasetypes.CreateProductUsecaseInput) (*model.Product, error) {
	validatedInput, err := createProductUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := createProductUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createProductUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/createProductUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := createProductUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.MemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              createProductUcase.createproductAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createProductUsecase",
			err,
		)
	}
	if accMemberAccess.Access.ProductAccesses.ProductApproval != nil {
		if *accMemberAccess.Access.ProductAccesses.ProductApproval {
			validatedInput.CreateProduct.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	productToCreate := &model.InternalCreateProduct{}
	jsonTemp, _ := json.Marshal(validatedInput.CreateProduct)
	json.Unmarshal(jsonTemp, productToCreate)

	for i, descriptivePhoto := range validatedInput.CreateProduct.Photos {
		productToCreate.Photos[i].Photo.File = descriptivePhoto.Photo.File
	}

	productToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdProduct, err := createProductUcase.createProductRepo.RunTransaction(
		productToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/createProductUsecase",
			err,
		)
	}

	return createdProduct, nil
}
