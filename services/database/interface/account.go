package databaseserviceinterface

import (
	model "github.com/horeekaa/backend/model"
)

type AccountService interface {
	FindByID(ID interface{}, serviceOptions *ServiceOptions) (chan<- *model.Account, chan<- error)
	FindOne(query map[string]interface{}, serviceOptions *ServiceOptions) (chan<- *model.Account, chan<- error)
	Find(query map[string]interface{}, serviceOptions *ServiceOptions) (chan<- []*model.Account, chan<- error)
	Create(input *model.CreateAccount, serviceOptions *ServiceOptions) (chan<- *model.Account, chan<- error)
	Update(ID interface{}, updateData *model.UpdateAccount, serviceOptions *ServiceOptions) (chan<- *model.Account, chan<- error)
}
