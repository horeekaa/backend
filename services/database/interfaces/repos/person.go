package databaseserviceinterface

import (
	model "github.com/horeekaa/backend/model"
)

type PersonService interface {
	FindByID(ID interface{}, serviceOptions *ServiceOptions) (chan<- *model.Person, chan<- error)
	FindOne(query map[string]interface{}, serviceOptions *ServiceOptions) (chan<- *model.Person, chan<- error)
	Find(query map[string]interface{}, serviceOptions *ServiceOptions) (chan<- []*model.Person, chan<- error)
	Create(input *model.CreatePerson, serviceOptions *ServiceOptions) (chan<- *model.Person, chan<- error)
	Update(ID interface{}, updateData *model.UpdatePerson, serviceOptions *ServiceOptions) (chan<- *model.Person, chan<- error)
}
