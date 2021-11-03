package purchaseordertosupplydomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type processPurchaseOrderToSupplyTransactionComponent struct {
	addressDataSource               databaseaddressdatasourceinterfaces.AddressDataSource
	tagDataSource                   databasetagdatasourceinterfaces.TagDataSource
	taggingDataSource               databasetaggingdatasourceinterfaces.TaggingDataSource
	memberAccessDataSource          databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
}

func NewProcessPurchaseOrderToSupplyTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
) (purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyTransactionComponent, error) {
	return &processPurchaseOrderToSupplyTransactionComponent{
		addressDataSource:               addressDataSource,
		tagDataSource:                   tagDataSource,
		taggingDataSource:               taggingDataSource,
		memberAccessDataSource:          memberAccessDataSource,
		purchaseOrderToSupplyDataSource: purchaseOrderToSupplyDataSource,
	}, nil
}

func (processPOToSupplyTrx *processPurchaseOrderToSupplyTransactionComponent) PreTransaction(
	input *model.PurchaseOrderToSupply,
) (*model.PurchaseOrderToSupply, error) {
	return input, nil
}

func (processPOToSupplyTrx *processPurchaseOrderToSupplyTransactionComponent) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input *model.PurchaseOrderToSupply,
) ([]*model.InternalCreateNotification, error) {
	taggings, err := processPOToSupplyTrx.taggingDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"tag._id": map[string]interface{}{
				"$in": funk.Map(
					input.Tags,
					func(t *model.TagForPurchaseOrderToSupply) interface{} {
						return t.ID
					},
				),
			},
			"taggingType": model.TaggingTypeOrganization,
			"isActive":    true,
		},
		&mongodbcoretypes.PaginationOptions{},
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/processPurchaseOrderToSupplyTransactionComponent",
			err,
		)
	}

	notifsToCreate := []*model.InternalCreateNotification{}

	for _, tagging := range taggings {
		address, err := processPOToSupplyTrx.addressDataSource.GetMongoDataSource().FindOne(
			map[string]interface{}{
				"object._id":             tagging.Organization.ID,
				"addressRegionGroup._id": input.AddressRegionGroup.ID,
			},
			operationOption,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/processPurchaseOrderToSupplyTransactionComponent",
				err,
			)
		}
		if address == nil {
			continue
		}

		loadedMemberAccessesChan := make(chan bool)
		loadedTagChan := make(chan bool)
		errChan := make(chan error)

		memberAccesses := []*model.MemberAccess{}
		go func() {
			loadedMemberAccesses, err := processPOToSupplyTrx.memberAccessDataSource.GetMongoDataSource().Find(
				map[string]interface{}{
					"organization._id":  tagging.Organization.ID,
					"organization.type": model.OrganizationTypePartner,
					"status":            model.MemberAccessStatusActive,
				},
				&mongodbcoretypes.PaginationOptions{},
				operationOption,
			)
			if err != nil {
				errChan <- horeekaacoreexceptiontofailure.ConvertException(
					"/processPurchaseOrderToSupplyTransactionComponent",
					err,
				)
				return
			}
			memberAccesses = append(memberAccesses, loadedMemberAccesses...)
			loadedMemberAccessesChan <- true
		}()

		tag := &model.Tag{}
		go func() {
			loadedTag, err := processPOToSupplyTrx.tagDataSource.GetMongoDataSource().FindByID(
				tagging.Tag.ID,
				operationOption,
			)
			if err != nil {
				errChan <- horeekaacoreexceptiontofailure.ConvertException(
					"/processPurchaseOrderToSupplyTransactionComponent",
					err,
				)
				return
			}
			tag = loadedTag
			loadedTagChan <- true
		}()

		for i := 0; i < 2; {
			select {
			case err := <-errChan:
				return nil, err
			case _ = <-loadedMemberAccessesChan:
				i++
			case _ = <-loadedTagChan:
				i++
			}
		}

		for _, memberAccess := range memberAccesses {
			notifToCreate := &model.InternalCreateNotification{
				NotificationCategory: model.NotificationCategoryPurchaseOrderSupplyBroadcast,
				RecipientAccount: &model.ObjectIDOnly{
					ID: &memberAccess.Account.ID,
				},
				PayloadOptions: &model.PayloadOptionsInput{
					PurchaseOrderToSupplyBroadcastPayload: &model.PurchaseOrderToSupplyBroadcastPayloadInput{
						BroadCastedByTag:      &model.TagForNotifPayloadInput{},
						PurchaseOrderToSupply: &model.PurchaseOrderToSupplyForNotifPayloadInput{},
					},
				},
			}
			jsonTag, _ := json.Marshal(tag)
			json.Unmarshal(jsonTag, &notifToCreate.PayloadOptions.PurchaseOrderToSupplyBroadcastPayload.BroadCastedByTag)

			jsonPOToSupply, _ := json.Marshal(input)
			json.Unmarshal(jsonPOToSupply, &notifToCreate.PayloadOptions.PurchaseOrderToSupplyBroadcastPayload.PurchaseOrderToSupply)

			notifsToCreate = append(notifsToCreate, notifToCreate)
		}
	}

	_, err = processPOToSupplyTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": input.ID,
		},
		&model.DatabaseUpdatePurchaseOrderToSupply{
			Status: func(m model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus {
				return &m
			}(model.PurchaseOrderToSupplyStatusProcessed),
		},
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/processPurchaseOrderToSupplyTransactionComponent",
			err,
		)
	}

	return notifsToCreate, nil
}
