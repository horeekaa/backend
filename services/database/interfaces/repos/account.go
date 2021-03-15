package databaseservicerepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type AccountService interface {
	FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error)
	FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error)
	Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.Account, chan error)
	Create(input *model.CreateAccount, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error)
	Update(ID interface{}, updateData *model.UpdateAccount, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error)
}
