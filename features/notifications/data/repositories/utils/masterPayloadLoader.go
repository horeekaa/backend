package notificationdomainrepositoryutilities

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
	"github.com/horeekaa/backend/model"
)

type masterPayloadLoader struct {
	invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader
}

func NewMasterPayloadLoader(
	invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader,
) (notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader, error) {
	return &masterPayloadLoader{
		invitationPayloadLoader: invitationPayloadLoader,
	}, nil
}

func (masterPayload *masterPayloadLoader) TransactionBody(
	operationOptions *mongodbcoretypes.OperationOptions,
	notification *model.InternalCreateNotification,
) (bool, error) {
	_, err := masterPayload.invitationPayloadLoader.TransactionBody(
		operationOptions,
		notification,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
