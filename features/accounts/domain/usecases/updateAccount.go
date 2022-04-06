package accountpresentationusecases

import (
	"encoding/json"

	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updateAccountUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	updateAccountRepo                accountdomainrepositoryinterfaces.UpdateAccountRepository
	updateOwnedaccountAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                     string
}

func NewUpdateAccountUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	updateAccountRepo accountdomainrepositoryinterfaces.UpdateAccountRepository,
) (accountpresentationusecaseinterfaces.UpdateAccountUsecase, error) {
	return &updateAccountUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		updateAccountRepo,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountUpdateOwned: func(b bool) *bool { return &b }(true),
			},
		},
		"UpdateAccountUsecase",
	}, nil
}

func (updateAccountUcase *updateAccountUsecase) validation(input accountpresentationusecasetypes.UpdateAccountUsecaseInput) (accountpresentationusecasetypes.UpdateAccountUsecaseInput, error) {
	if &input.Context == nil {
		return accountpresentationusecasetypes.UpdateAccountUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				updateAccountUcase.pathIdentity,
				nil,
			)
	}

	return input, nil
}

func (updateAccountUcase *updateAccountUsecase) Execute(input accountpresentationusecasetypes.UpdateAccountUsecaseInput) (*model.Account, error) {
	validatedInput, err := updateAccountUcase.validation(input)
	if err != nil {
		return nil, err
	}
	accountToUpdate := &model.InternalUpdateAccount{}
	jsonTemp, _ := json.Marshal(validatedInput.UpdateAccount)
	json.Unmarshal(jsonTemp, accountToUpdate)

	account, err := updateAccountUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAccountUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			updateAccountUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrganizationBased := model.MemberAccessRefTypeOrganizationsBased
	accMemberAccess, err := updateAccountUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrganizationBased,
				Status: func(s model.MemberAccessStatus) *model.MemberAccessStatus {
					return &s
				}(model.MemberAccessStatusActive),
				ProposalStatus: func(e model.EntityProposalStatus) *model.EntityProposalStatus {
					return &e
				}(model.EntityProposalStatusApproved),
				InvitationAccepted: func(b bool) *bool {
					return &b
				}(true),
			},
			QueryMode: true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAccountUcase.pathIdentity,
			err,
		)
	}

	// if update across accounts is not allowed, check access for update owned account
	if accessible := funk.GetOrElse(
		funk.Get(accMemberAccess, "Access.AccountAccesses.AccountUpdate"), false,
	).(bool); !accessible {
		if accessible := funk.GetOrElse(
			funk.Get(accMemberAccess, "Access.AccountAccesses.AccountUpdateOwned"), false,
		).(bool); accessible {
			if account.ID.Hex() != accountToUpdate.ID.Hex() {
				return nil, horeekaacoreerror.NewErrorObject(
					horeekaacorefailureenums.FeatureNotAccessibleByAccount,
					updateAccountUcase.pathIdentity,
					nil,
				)
			}
		} else {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacorefailureenums.FeatureNotAccessibleByAccount,
				updateAccountUcase.pathIdentity,
				nil,
			)
		}
	}

	updateAccountOutput, err := updateAccountUcase.updateAccountRepo.RunTransaction(
		accountToUpdate,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			updateAccountUcase.pathIdentity,
			err,
		)
	}

	return updateAccountOutput, nil
}
