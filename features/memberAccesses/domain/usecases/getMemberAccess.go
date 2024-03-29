package memberaccesspresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getMemberAccessUsecase struct {
	getMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	pathIdentity              string
}

func NewGetMemberAccessUsecase(
	getMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
) (memberaccesspresentationusecaseinterfaces.GetMemberAccessUsecase, error) {
	return &getMemberAccessUsecase{
		getMemberAccessRepository,
		"GetMemberAccessUsecase",
	}, nil
}

func (getMmbAccUcase *getMemberAccessUsecase) validation(
	input *model.InternalMemberAccessFilterFields,
) (*model.InternalMemberAccessFilterFields, error) {
	return input, nil
}

func (getMmbAccUcase *getMemberAccessUsecase) Execute(
	filterFields *model.InternalMemberAccessFilterFields,
) (*model.MemberAccess, error) {
	validatedFilterFields, err := getMmbAccUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	memberAccess, err := getMmbAccUcase.getMemberAccessRepository.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: validatedFilterFields,
			QueryMode:                true,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getMmbAccUcase.pathIdentity,
			err,
		)
	}
	return memberAccess, nil
}
