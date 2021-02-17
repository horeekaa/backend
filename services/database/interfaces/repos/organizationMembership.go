package databaseservicerepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type OrganizationMembershipService interface {
	FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error)
	FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error)
	Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- []*model.OrganizationMembership, chan<- error)
	Create(input *model.CreateOrganizationMembership, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error)
	Update(ID interface{}, updateData *model.UpdateOrganizationMembership, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.OrganizationMembership, chan<- error)
}
