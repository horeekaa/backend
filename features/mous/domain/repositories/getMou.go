package moudomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetMouRepository interface {
	Execute(filterFields *model.MouFilterFields) (*model.Mou, error)
}
