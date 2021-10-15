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
}

func NewCreateAddressRegionGroupTransactionComponent(
	addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent, error) {
	return &createAddressRegionGroupTransactionComponent{
		addressRegionGroupDataSource: addressRegionGroupDataSource,
		loggingDataSource:            loggingDataSource,
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
	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createAddressRegionGroupTrx.GetCurrentObjectID()

	loc, _ := time.LoadLocation("Asia/Bangkok")
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	input.PublicID = func(s ...string) *string { joinedString := strings.Join(s, "/"); return &joinedString }(
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
				ID: input.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *input.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createAddressRegionGroup",
			err,
		)
	}

	input.ID = generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	addressRegionGroupToCreate := &model.DatabaseCreateAddressRegionGroup{}

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, addressRegionGroupToCreate)
	json.Unmarshal(jsonTemp, &addressRegionGroupToCreate.ProposedChanges)

	newAddressRegionGroup, err := createAddressRegionGroupTrx.addressRegionGroupDataSource.GetMongoDataSource().Create(
		addressRegionGroupToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createAddressRegionGroup",
			err,
		)
	}
	createAddressRegionGroupTrx.generatedObjectID = nil

	return newAddressRegionGroup, nil
}
