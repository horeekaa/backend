package databaseservicerepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type MemberAccessService interface {
	FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error)
	FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error)
	Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.MemberAccess, chan error)
	Create(input *model.CreateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error)
	Update(ID interface{}, updateData *model.UpdateMemberAccess, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.MemberAccess, chan error)
}
