package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	productvariantpresentationusecaseinterfaces "github.com/horeekaa/backend/features/productVariants/presentation/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, createProduct model.CreateProduct) (*model.Product, error) {
	var createProductUsecase productpresentationusecaseinterfaces.CreateProductUsecase
	container.Make(&createProductUsecase)
	return createProductUsecase.Execute(
		productpresentationusecasetypes.CreateProductUsecaseInput{
			Context:       ctx,
			CreateProduct: &createProduct,
		},
	)
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, updateProduct model.UpdateProduct) (*model.Product, error) {
	var updateProductUsecase productpresentationusecaseinterfaces.UpdateProductUsecase
	container.Make(&updateProductUsecase)
	return updateProductUsecase.Execute(
		productpresentationusecasetypes.UpdateProductUsecaseInput{
			Context:       ctx,
			UpdateProduct: &updateProduct,
		},
	)
}

func (r *productResolver) Photos(ctx context.Context, obj *model.Product) ([]*model.DescriptivePhoto, error) {
	var getDescriptivePhotoUsecase descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase
	container.Make(&getDescriptivePhotoUsecase)

	descriptivePhotos := []*model.DescriptivePhoto{}
	if obj.Photos != nil {
		for _, photo := range obj.Photos {
			descriptivePhoto, err := getDescriptivePhotoUsecase.Execute(
				&model.DescriptivePhotoFilterFields{
					ID: &photo.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			descriptivePhotos = append(descriptivePhotos, descriptivePhoto)
		}
	}
	return descriptivePhotos, nil
}

func (r *productResolver) Variants(ctx context.Context, obj *model.Product) ([]*model.ProductVariant, error) {
	var getProductVariantUsecase productvariantpresentationusecaseinterfaces.GetProductVariantUsecase
	container.Make(&getProductVariantUsecase)

	productVariants := []*model.ProductVariant{}
	if obj.Variants != nil {
		for _, variant := range obj.Variants {
			productVariant, err := getProductVariantUsecase.Execute(
				&model.ProductVariantFilterFields{
					ID: &variant.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			productVariants = append(productVariants, productVariant)
		}
	}
	return productVariants, nil
}

func (r *productResolver) Taggings(ctx context.Context, obj *model.Product) ([]*model.Tagging, error) {
	var getTaggingUsecase taggingpresentationusecaseinterfaces.GetTaggingUsecase
	container.Make(&getTaggingUsecase)

	taggings := []*model.Tagging{}
	if obj.Taggings != nil {
		for _, tagg := range obj.Taggings {
			tagging, err := getTaggingUsecase.Execute(
				&model.TaggingFilterFields{
					ID: &tagg.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			taggings = append(taggings, tagging)
		}
	}
	return taggings, nil
}

func (r *productResolver) SubmittingAccount(ctx context.Context, obj *model.Product) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: &model.AccountFilterFields{
				ID: &obj.SubmittingAccount.ID,
			},
		},
	)
}

func (r *productResolver) RecentApprovingAccount(ctx context.Context, obj *model.Product) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)

	var filterFields *model.AccountFilterFields
	if obj.RecentApprovingAccount != nil {
		filterFields = &model.AccountFilterFields{}
		filterFields.ID = &obj.RecentApprovingAccount.ID
	}
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: filterFields,
		},
	)
}

func (r *productResolver) RecentLog(ctx context.Context, obj *model.Product) (*model.Logging, error) {
	var getLoggingUsecase loggingpresentationusecaseinterfaces.GetLoggingUsecase
	container.Make(&getLoggingUsecase)

	var filterFields *model.LoggingFilterFields
	if obj.RecentLog != nil {
		filterFields = &model.LoggingFilterFields{}
		filterFields.ID = &obj.RecentLog.ID
	}
	return getLoggingUsecase.Execute(
		filterFields,
	)
}

func (r *queryResolver) Products(ctx context.Context, filterFields model.ProductFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Product, error) {
	var getProductsUsecase productpresentationusecaseinterfaces.GetAllProductUsecase
	container.Make(&getProductsUsecase)
	return getProductsUsecase.Execute(
		productpresentationusecasetypes.GetAllProductUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
