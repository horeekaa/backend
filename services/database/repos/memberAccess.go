package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	databaseinstancereferences "github.com/horeekaa/backend/repositories/databaseClient/instanceReferences/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type memberAccessService struct {
	memberAccessRepo *databaseinstancereferences.MemberAccessRepo
}

func NewMemberAccessService(memberAccessRepo databaseinstancereferences.MemberAccessRepo) (databaseservicerepointerfaces.MemberAccessService, error) {
	return &memberAccessService{
		&memberAccessRepo,
	}, nil
}

func (memberAccessSvc *memberAccessService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error) {
	memberAccessChn := make(chan *model.MemberAccess)
	errorChn := make(chan error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo.Instance).FindByID(ID, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessService/FindByID",
				&err,
			)
			return
		}

		memberAccessChn <- memberAccess
	}()

	return memberAccessChn, errorChn
}

func (memberAccessSvc *memberAccessService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error) {
	memberAccessChn := make(chan *model.MemberAccess)
	errorChn := make(chan error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo.Instance).FindOne(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessService/FindOne",
				&err,
			)
			return
		}

		memberAccessChn <- memberAccess
	}()

	return memberAccessChn, errorChn
}

func (memberAccessSvc *memberAccessService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.MemberAccess, chan error) {
	memberAccessesChn := make(chan []*model.MemberAccess)
	errorChn := make(chan error)

	go func() {
		memberAccesses, err := (*memberAccessSvc.memberAccessRepo.Instance).Find(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessService/Find",
				&err,
			)
			return
		}

		memberAccessesChn <- memberAccesses
	}()

	return memberAccessesChn, errorChn
}

func (memberAccessSvc *memberAccessService) Create(input *model.CreateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error) {
	memberAccessChn := make(chan *model.MemberAccess)
	errorChn := make(chan error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo.Instance).Create(input, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessService/Create",
				&err,
			)
			return
		}

		memberAccessChn <- memberAccess
	}()

	return memberAccessChn, errorChn
}

func (memberAccessSvc *memberAccessService) Update(ID interface{}, updateData *model.UpdateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error) {
	memberAccessChn := make(chan *model.MemberAccess)
	errorChn := make(chan error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo.Instance).Update(ID, updateData, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessService/Update",
				&err,
			)
			return
		}

		memberAccessChn <- memberAccess
	}()

	return memberAccessChn, errorChn
}
