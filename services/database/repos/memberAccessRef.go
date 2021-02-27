package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type memberAccessRefService struct {
	memberAccessRefRepo *mongorepointerfaces.MemberAccessRefRepoMongo
}

func NewMemberAccessRefService(memberAccessRefRepo mongorepointerfaces.MemberAccessRefRepoMongo) (databaseservicerepointerfaces.MemberAccessRefService, error) {
	return &memberAccessRefService{
		&memberAccessRefRepo,
	}, nil
}

func (memberAccessRefSvc *memberAccessRefService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error) {
	memberAccessRefChn := make(chan *model.MemberAccessRef)
	errorChn := make(chan error)

	go func() {
		memberAccessRef, err := (*memberAccessRefSvc.memberAccessRefRepo).FindByID(ID, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessRefService/FindByID",
				&err,
			)
			return
		}

		memberAccessRefChn <- memberAccessRef
	}()

	return memberAccessRefChn, errorChn
}

func (memberAccessRefSvc *memberAccessRefService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error) {
	memberAccessRefChn := make(chan *model.MemberAccessRef)
	errorChn := make(chan error)

	go func() {
		memberAccessRef, err := (*memberAccessRefSvc.memberAccessRefRepo).FindOne(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessRefService/FindOne",
				&err,
			)
			return
		}

		memberAccessRefChn <- memberAccessRef
	}()

	return memberAccessRefChn, errorChn
}

func (memberAccessRefSvc *memberAccessRefService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.MemberAccessRef, chan error) {
	memberAccessRefsChn := make(chan []*model.MemberAccessRef)
	errorChn := make(chan error)

	go func() {
		memberAccessRefs, err := (*memberAccessRefSvc.memberAccessRefRepo).Find(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessRefService/Find",
				&err,
			)
			return
		}

		memberAccessRefsChn <- memberAccessRefs
	}()

	return memberAccessRefsChn, errorChn
}

func (memberAccessRefSvc *memberAccessRefService) Create(input *model.CreateMemberAccessRef, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error) {
	memberAccessRefChn := make(chan *model.MemberAccessRef)
	errorChn := make(chan error)

	go func() {
		memberAccessRef, err := (*memberAccessRefSvc.memberAccessRefRepo).Create(input, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessRefService/Create",
				&err,
			)
			return
		}

		memberAccessRefChn <- memberAccessRef
	}()

	return memberAccessRefChn, errorChn
}

func (memberAccessRefSvc *memberAccessRefService) Update(ID interface{}, updateData *model.UpdateMemberAccessRef, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error) {
	memberAccessRefChn := make(chan *model.MemberAccessRef)
	errorChn := make(chan error)

	go func() {
		memberAccessRef, err := (*memberAccessRefSvc.memberAccessRefRepo).Update(ID, updateData, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/memberAccessRefService/Update",
				&err,
			)
			return
		}

		memberAccessRefChn <- memberAccessRef
	}()

	return memberAccessRefChn, errorChn
}
