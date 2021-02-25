package databaseservicerepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type MemberAccessRefService interface {
	FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error)
	FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error)
	Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.MemberAccessRef, chan error)
	Create(input *model.CreateMemberAccessRef, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error)
	Update(ID interface{}, updateData *model.UpdateMemberAccessRef, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccessRef, chan error)
}
