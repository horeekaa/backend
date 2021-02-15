package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type memberAccessService struct {
	memberAccessRepo *mongorepointerfaces.MemberAccessRepoMongo
}

func NewMemberAccessService(memberAccessRepo *mongorepointerfaces.MemberAccessRepoMongo) (databaseservicerepointerfaces.MemberAccessService, error) {
	return &memberAccessService{
		memberAccessRepo,
	}, nil
}

func (memberAccessSvc *memberAccessService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.MemberAccess, chan<- error) {
	memberAccessChn := make(chan<- *model.MemberAccess)
	errorChn := make(chan<- error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo).FindByID(ID, (*serviceOptions).OperationOptions)
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

func (memberAccessSvc *memberAccessService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.MemberAccess, chan<- error) {
	memberAccessChn := make(chan<- *model.MemberAccess)
	errorChn := make(chan<- error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo).FindOne(query, (*serviceOptions).OperationOptions)
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

func (memberAccessSvc *memberAccessService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- []*model.MemberAccess, chan<- error) {
	memberAccessesChn := make(chan<- []*model.MemberAccess)
	errorChn := make(chan<- error)

	go func() {
		memberAccesses, err := (*memberAccessSvc.memberAccessRepo).Find(query, (*serviceOptions).OperationOptions)
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

func (memberAccessSvc *memberAccessService) Create(input *model.CreateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.MemberAccess, chan<- error) {
	memberAccessChn := make(chan<- *model.MemberAccess)
	errorChn := make(chan<- error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo).Create(input, (*serviceOptions).OperationOptions)
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

func (memberAccessSvc *memberAccessService) Update(ID interface{}, updateData *model.UpdateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.MemberAccess, chan<- error) {
	memberAccessChn := make(chan<- *model.MemberAccess)
	errorChn := make(chan<- error)

	go func() {
		memberAccess, err := (*memberAccessSvc.memberAccessRepo).Update(ID, updateData, (*serviceOptions).OperationOptions)
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
