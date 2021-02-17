package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type organizationMembershipService struct {
	organizationMembershipRepo *mongorepointerfaces.OrganizationMembershipRepoMongo
}

func NewOrganizationMembershipService(organizationMembershipRepo *mongorepointerfaces.OrganizationMembershipRepoMongo) (databaseservicerepointerfaces.OrganizationMembershipService, error) {
	return &organizationMembershipService{
		organizationMembershipRepo,
	}, nil
}

func (organizationMembershipSvc *organizationMembershipService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error) {
	organizationMembershipChn := make(chan<- *model.OrganizationMembership)
	errorChn := make(chan<- error)

	go func() {
		organizationMembership, err := (*organizationMembershipSvc.organizationMembershipRepo).FindByID(ID, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/organizationMembershipService/FindByID",
				&err,
			)
			return
		}

		organizationMembershipChn <- organizationMembership
	}()

	return organizationMembershipChn, errorChn
}

func (organizationMembershipSvc *organizationMembershipService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error) {
	organizationMembershipChn := make(chan<- *model.OrganizationMembership)
	errorChn := make(chan<- error)

	go func() {
		organizationMembership, err := (*organizationMembershipSvc.organizationMembershipRepo).FindOne(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/organizationMembershipService/FindOne",
				&err,
			)
			return
		}

		organizationMembershipChn <- organizationMembership
	}()

	return organizationMembershipChn, errorChn
}

func (organizationMembershipSvc *organizationMembershipService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- []*model.OrganizationMembership, chan<- error) {
	organizationMembershipesChn := make(chan<- []*model.OrganizationMembership)
	errorChn := make(chan<- error)

	go func() {
		organizationMembershipes, err := (*organizationMembershipSvc.organizationMembershipRepo).Find(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/organizationMembershipService/Find",
				&err,
			)
			return
		}

		organizationMembershipesChn <- organizationMembershipes
	}()

	return organizationMembershipesChn, errorChn
}

func (organizationMembershipSvc *organizationMembershipService) Create(input *model.CreateOrganizationMembership, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error) {
	organizationMembershipChn := make(chan<- *model.OrganizationMembership)
	errorChn := make(chan<- error)

	go func() {
		organizationMembership, err := (*organizationMembershipSvc.organizationMembershipRepo).Create(input, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/organizationMembershipService/Create",
				&err,
			)
			return
		}

		organizationMembershipChn <- organizationMembership
	}()

	return organizationMembershipChn, errorChn
}

func (organizationMembershipSvc *organizationMembershipService) Update(ID interface{}, updateData *model.UpdateOrganizationMembership, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error) {
	organizationMembershipChn := make(chan<- *model.OrganizationMembership)
	errorChn := make(chan<- error)

	go func() {
		organizationMembership, err := (*organizationMembershipSvc.organizationMembershipRepo).Update(ID, updateData, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/organizationMembershipService/Update",
				&err,
			)
			return
		}

		organizationMembershipChn <- organizationMembership
	}()

	return organizationMembershipChn, errorChn
}
