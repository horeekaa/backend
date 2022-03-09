package addressdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createAddressTransactionComponent struct {
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource
	addressLoader     addressdomainrepositoryutilityinterfaces.AddressLoader
	generatedObjectID *primitive.ObjectID
	pathIdentity      string
}

func NewCreateAddressTransactionComponent(
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
) (addressdomainrepositoryinterfaces.CreateAddressTransactionComponent, error) {
	return &createAddressTransactionComponent{
		addressDataSource: addressDataSource,
		loggingDataSource: loggingDataSource,
		addressLoader:     addressLoader,
		pathIdentity:      "CreateAddressComponent",
	}, nil
}

func (createAddressTrx *createAddressTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createAddressTrx.addressDataSource.GetMongoDataSource().GenerateObjectID()
	createAddressTrx.generatedObjectID = &generatedObjectID
	return *createAddressTrx.generatedObjectID
}

func (createAddressTrx *createAddressTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createAddressTrx.generatedObjectID == nil {
		generatedObjectID := createAddressTrx.addressDataSource.GetMongoDataSource().GenerateObjectID()
		createAddressTrx.generatedObjectID = &generatedObjectID
	}
	return *createAddressTrx.generatedObjectID
}

func (createAddrTrx *createAddressTransactionComponent) PreTransaction(
	createaddressInput *model.InternalCreateAddress,
) (*model.InternalCreateAddress, error) {
	return createaddressInput, nil
}

func (createAddrTrx *createAddressTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateAddress,
) (*model.Address, error) {
	addressToCreate := &model.DatabaseCreateAddress{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, addressToCreate)
	addressToCreate.ResolvedGeocoding = &model.ResolvedGeocodingInput{}
	addressToCreate.AddressRegionGroup = &model.AddressRegionGroupForAddressInput{}
	_, err := createAddrTrx.addressLoader.Execute(
		session,
		&addressdomainrepositorytypes.LatLngGeocode{
			Latitude:  addressToCreate.Latitude,
			Longitude: addressToCreate.Longitude,
		},
		addressToCreate.ResolvedGeocoding,
		addressToCreate.AddressRegionGroup,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAddrTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*addressToCreate)
	generatedObjectID := createAddrTrx.GetCurrentObjectID()
	loggingOutput, err := createAddrTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Address",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: addressToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *addressToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAddrTrx.pathIdentity,
			err,
		)
	}

	addressToCreate.ID = generatedObjectID
	addressToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *addressToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		addressToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: addressToCreate.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(addressToCreate)
	json.Unmarshal(jsonTemp, &addressToCreate.ProposedChanges)

	createdAddress, err := createAddrTrx.addressDataSource.GetMongoDataSource().Create(
		addressToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAddrTrx.pathIdentity,
			err,
		)
	}

	return createdAddress, nil
}
