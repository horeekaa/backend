package addressregiongroupdomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createAddressRegionGroupTransactionComponent struct {
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
	loggingDataSource            databaseloggingdatasourceinterfaces.LoggingDataSource
	generatedObjectID            *primitive.ObjectID
	pathIdentity                 string
}

func NewCreateAddressRegionGroupTransactionComponent(
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent, error) {
	return &createAddressRegionGroupTransactionComponent{
		addressRegionGroupDataSource: addressRegionGroupDataSource,
		loggingDataSource:            loggingDataSource,
		pathIdentity:                 "CreateAddressRegionGroupComponent",
	}, nil
}

func (createAddressRegionGroupTrx *createAddressRegionGroupTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().GenerateObjectID()
	createAddressRegionGroupTrx.generatedObjectID = &generatedObjectID
	return *createAddressRegionGroupTrx.generatedObjectID
}

func (createAddressRegionGroupTrx *createAddressRegionGroupTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createAddressRegionGroupTrx.generatedObjectID == nil {
		generatedObjectID := createAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().GenerateObjectID()
		createAddressRegionGroupTrx.generatedObjectID = &generatedObjectID
	}
	return *createAddressRegionGroupTrx.generatedObjectID
}

func (createAddressRegionGroupTrx *createAddressRegionGroupTransactionComponent) PreTransaction(
	input *model.InternalCreateAddressRegionGroup,
) (*model.InternalCreateAddressRegionGroup, error) {
	return input, nil
}

func (createAddressRegionGroupTrx *createAddressRegionGroupTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	addressRegionGroupToCreate := &model.DatabaseCreateAddressRegionGroup{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, addressRegionGroupToCreate)

	newDocumentJson, _ := json.Marshal(*addressRegionGroupToCreate)
	generatedObjectID := createAddressRegionGroupTrx.GetCurrentObjectID()
	loc, _ := time.LoadLocation("Asia/Bangkok")
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	addressRegionGroupToCreate.PublicID = func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
		"ARG",
		time.Now().In(loc).Format("20060102"),
		strings.ToUpper(
			strings.Join(
				splittedId[len(splittedId)-4:],
				"",
			),
		),
	)
	loggingOutput, err := createAddressRegionGroupTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "AddressRegionGroup",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: addressRegionGroupToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *addressRegionGroupToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAddressRegionGroupTrx.pathIdentity,
			err,
		)
	}

	addressRegionGroupToCreate.ID = generatedObjectID
	addressRegionGroupToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *addressRegionGroupToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		addressRegionGroupToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: addressRegionGroupToCreate.SubmittingAccount.ID}
	}
	currentTime := time.Now()
	addressRegionGroupToCreate.CreatedAt = &currentTime
	addressRegionGroupToCreate.UpdatedAt = &currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if addressRegionGroupToCreate.ProposalStatus == nil {
		addressRegionGroupToCreate.ProposalStatus = &defaultProposalStatus
	}

	jsonTemp, _ = json.Marshal(addressRegionGroupToCreate)
	json.Unmarshal(jsonTemp, &addressRegionGroupToCreate.ProposedChanges)

	newAddressRegionGroup, err := createAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().Create(
		addressRegionGroupToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAddressRegionGroupTrx.pathIdentity,
			err,
		)
	}
	createAddressRegionGroupTrx.generatedObjectID = nil

	return newAddressRegionGroup, nil
}
