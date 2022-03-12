package taggingpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type bulkCreateTaggingUsecase struct {
	getAccountFromAuthDataRepo      accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo      memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	bulkCreateTaggingRepo           taggingdomainrepositoryinterfaces.BulkCreateTaggingRepository
	bulkCreateTaggingAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                    string
}

func NewBulkCreateTaggingUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	bulkCreateTaggingRepo taggingdomainrepositoryinterfaces.BulkCreateTaggingRepository,
) (taggingpresentationusecaseinterfaces.BulkCreateTaggingUsecase, error) {
	return &bulkCreateTaggingUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		bulkCreateTaggingRepo,
		&model.MemberAccessRefOptionsInput{
			BulkTaggingAccesses: &model.BulkTaggingAccessesInput{
				BulkTaggingCreate: func(b bool) *bool { return &b }(true),
			},
		},
		"BulkCreateTaggingUsecase",
	}, nil
}

func (bulkCreateTaggingUcase *bulkCreateTaggingUsecase) validation(input taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput) (taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput, error) {
	if &input.Context == nil {
		return taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				bulkCreateTaggingUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (bulkCreateTaggingUcase *bulkCreateTaggingUsecase) Execute(input taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput) ([]*model.Tagging, error) {
	validatedInput, err := bulkCreateTaggingUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := bulkCreateTaggingUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			bulkCreateTaggingUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			bulkCreateTaggingUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := bulkCreateTaggingUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              bulkCreateTaggingUcase.bulkCreateTaggingAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			bulkCreateTaggingUcase.pathIdentity,
			err,
		)
	}
	if accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval != nil {
		if *accMemberAccess.Access.BulkTaggingAccesses.BulkTaggingApproval {
			validatedInput.BulkCreateTagging.ProposalStatus =
				func(i model.EntityProposalStatus) *model.EntityProposalStatus { return &i }(model.EntityProposalStatusApproved)
		}
	}

	taggingToCreate := &model.InternalCreateTagging{}
	jsonTemp, _ := json.Marshal(validatedInput.BulkCreateTagging)
	json.Unmarshal(jsonTemp, taggingToCreate)

	taggingToCreate.SubmittingAccount = &model.ObjectIDOnly{ID: &account.ID}
	createdTaggings, err := bulkCreateTaggingUcase.bulkCreateTaggingRepo.RunTransaction(
		taggingToCreate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			bulkCreateTaggingUcase.pathIdentity,
			err,
		)
	}

	return createdTaggings, nil
}
