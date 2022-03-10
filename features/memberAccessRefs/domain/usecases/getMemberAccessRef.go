package memberaccessrefpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getMemberAccessRefUsecase struct {
	getMemberAccessRefRepository memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository
	pathIdentity                 string
}

func NewGetMemberAccessRefUsecase(
	getMemberAccessRefRepository memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository,
) (memberaccessrefpresentationusecaseinterfaces.GetMemberAccessRefUsecase, error) {
	return &getMemberAccessRefUsecase{
		getMemberAccessRefRepository,
		"GetMemberAccessRefUsecase",
	}, nil
}

func (getMmbAccRefUcase *getMemberAccessRefUsecase) validation(
	input *model.MemberAccessRefFilterFields,
) (*model.MemberAccessRefFilterFields, error) {
	return input, nil
}

func (getMmbAccRefUcase *getMemberAccessRefUsecase) Execute(
	filterFields *model.MemberAccessRefFilterFields,
) (*model.MemberAccessRef, error) {
	validatedFilterFields, err := getMmbAccRefUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	memberAccessRef, err := getMmbAccRefUcase.getMemberAccessRefRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getMmbAccRefUcase.pathIdentity,
			err,
		)
	}
	return memberAccessRef, nil
}
