package databaseservicerepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type PersonService interface {
	FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.Person, chan<- error)
	FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.Person, chan<- error)
	Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- []*model.Person, chan<- error)
	Create(input *model.CreatePerson, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.Person, chan<- error)
	Update(ID interface{}, updateData *model.UpdatePerson, serviceOptions *databaseserviceoperations.ServiceOptions) (chan<- *model.Person, chan<- error)
}
