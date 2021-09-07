package mouitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type updateMouItemDomainRepositories struct {
	mouItemDataSource   databasemouitemdatasourceinterfaces.MouItemDataSource
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader
}

func NewUpdateMouItemDomainRepository(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
) (mouitemdomainrepositoryinterfaces.UpdateMouItemTransactionComponent, error) {
	return &updateMouItemDomainRepositories{
		mouItemDataSource:   mouItemDataSource,
		agreedProductLoader: agreedProductLoader,
	}, nil
}

func (updateMouItemTrx *updateMouItemDomainRepositories) PreTransaction(
	input *model.InternalUpdateMouItem,
) (*model.InternalUpdateMouItem, error) {
	return input, nil
}

func (updateMouItemTrx *updateMouItemDomainRepositories) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMouItemInput *model.InternalUpdateMouItem,
) (*model.MouItem, error) {
	mouItemToUpdate := &model.DatabaseUpdateMouItem{}
	jsonTemp, _ := json.Marshal(updateMouItemInput)
	json.Unmarshal(jsonTemp, mouItemToUpdate)

	if mouItemToUpdate.Product != nil {
		updateMouItemTrx.agreedProductLoader.TransactionBody(
			session,
			mouItemToUpdate.Product,
			mouItemToUpdate.AgreedProduct,
		)
	}

	updatedMouItem, err := updateMouItemTrx.mouItemDataSource.GetMongoDataSource().Update(
		mouItemToUpdate.ID,
		mouItemToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMouItem",
			err,
		)
	}

	return updatedMouItem, nil
}
