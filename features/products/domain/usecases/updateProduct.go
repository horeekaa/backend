package productpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updateProductUsecase struct {
	getAccountFromAuthDataRepo  accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo  memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	proposeUpdateProductRepo    productdomainrepositoryinterfaces.ProposeUpdateProductRepository
	approveUpdateProductRepo    productdomainrepositoryinterfaces.ApproveUpdateProductRepository
	updateproductAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                string
}

func NewUpdateProductUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	proposeUpdateProductRepo productdomainrepositoryinterfaces.ProposeUpdateProductRepository,
	approveUpdateProductRepo productdomainrepositoryinterfaces.ApproveUpdateProductRepository,
) (productpresentationusecaseinterfaces.UpdateProductUsecase, error) {
	return &updateProductUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		proposeUpdateProductRepo,
		approveUpdateProductRepo,
		&model.MemberAccessRefOptionsInput{
			ProductAccesses: &model.ProductAccessesInput{
				ProductUpdate: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateProductUsecase",
	}, nil
}

func (updateProductUcase *updateProductUsecase) validation(input productpresentationusecasetypes.UpdateProductUsecaseInput) (productpresentationusecasetypes.UpdateProductUsecaseInput, error) {
	if &input.Context == nil {
		return productpresentationusecasetypes.UpdateProductUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateProductUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateProductUcase *updateProductUsecase) Execute(input productpresentationusecasetypes.UpdateProductUsecaseInput) (*model.Product, error) {
	validatedInput, err := updateProductUcase.validation(input)
	if err != nil {
		return nil, err
	}
	productToUpdate := &model.InternalUpdateProduct{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateProduct)
	json.Unmarshal(jsonTemp, productToUpdate)

	account, err := updateProductUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateProductUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateProductUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateProductUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              updateProductUcase.updateproductAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateProductUcase.pathIdentity,
			err,
		)
	}

	for i, descriptivePhoto := range validatedInput.UpdateProduct.Photos {
		if descriptivePhoto.Photo != nil {
			productToUpdate.Photos[i].Photo.File = descriptivePhoto.Photo.File
		}
	}

	for i, productVariant := range validatedInput.UpdateProduct.Variants {
		if funk.Get(productVariant, "Photo.Photo") != nil {
			productToUpdate.Variants[i].Photo.Photo.File = productVariant.Photo.Photo.File
		}
	}

	// if user is only going to approve proposal
	if productToUpdate.ProposalStatus != nil {
		if accMemberAccess.Access.ProductAccesses.ProductApproval == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateProductUcase.pathIdentity,
				nil,
			)
		}
		if !*accMemberAccess.Access.ProductAccesses.ProductApproval {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateProductUcase.pathIdentity,
				nil,
			)
		}

		productToUpdate.RecentApprovingAccount = &model.ObjectIDOnly{ID: &account.ID}
		updateproductOutput, err := updateProductUcase.approveUpdateProductRepo.RunTransaction(
			productToUpdate,
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				updateProductUcase.pathIdentity,
				err,
			)
		}

		return updateproductOutput, nil
	}

	productToUpdate.ProposalStatus =
		func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusProposed)
	if accMemberAccess.Access.ProductAccesses.ProductApproval != nil {
		if *accMemberAccess.Access.ProductAccesses.ProductApproval {
			productToUpdate.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	productToUpdate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	updateproductOutput, err := updateProductUcase.proposeUpdateProductRepo.RunTransaction(
		productToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateProductUcase.pathIdentity,
			err,
		)
	}

	return updateproductOutput, nil
}
