package memberaccessrefdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource                    databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	loggingDataSource                            databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                          coreutilityinterfaces.MapProcessorUtility
	structComparisonUtility                      coreutilityinterfaces.StructComparisonUtility
	proposeUpdateMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefUsecaseComponent
}

func NewProposeUpdateMemberAccessRefTransactionComponent(
	MemberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
) (memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent, error) {
	return &proposeUpdateMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: MemberAccessRefDataSource,
		loggingDataSource:         loggingDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		structComparisonUtility:   structComparisonUtility,
	}, nil
}

func (updateProdTrx *proposeUpdateMemberAccessRefTransactionComponent) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateProdTrx.proposeUpdateMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateProdTrx *proposeUpdateMemberAccessRefTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.InternalUpdateMemberAccessRef, error) {
	if updateProdTrx.proposeUpdateMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return updateProdTrx.proposeUpdateMemberAccessRefUsecaseComponent.Validation(input)
}

func (updateProdTrx *proposeUpdateMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccessRef *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	existingMemberAccessRef, err := updateProdTrx.memberAccessRefDataSource.GetMongoDataSource().FindByID(
		updateMemberAccessRef.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}
	fieldChanges := []*model.FieldChangeDataInput{}

	updateProdTrx.structComparisonUtility.SetComparisonFunc(
		func(tag interface{}, field1 interface{}, field2 interface{}, tagString *interface{}) {
			if field1 == field2 {
				return
			}
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field1).Kind().String(),
				OldValue: fmt.Sprint(field2),
				NewValue: fmt.Sprint(field1),
			})
			*tagString = ""
		},
	)
	updateProdTrx.structComparisonUtility.SetPreDeepComparisonFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	updateProdTrx.structComparisonUtility.CompareStructs(
		*updateMemberAccessRef,
		*existingMemberAccessRef,
		&tagString,
	)

	loggingOutput, err := updateProdTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccessRef",
			Document: &model.ObjectIDOnly{
				ID: &existingMemberAccessRef.ID,
			},
			FieldChanges: fieldChanges,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateMemberAccessRef.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateMemberAccessRef.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}
	updateMemberAccessRef.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateMemberAccessRef := &model.InternalUpdateMemberAccessRef{
		ID: updateMemberAccessRef.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccessRef)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	var updateMemberAccessRefMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jsonUpdate, &updateMemberAccessRefMap)

	updateProdTrx.mapProcessorUtility.RemoveNil(updateMemberAccessRefMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessRefMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	if updateMemberAccessRef.ProposalStatus != nil {
		if *updateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateMemberAccessRef.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateMemberAccessRef)
		}
	}

	updatedMemberAccessRef, err := updateProdTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
		fieldsToUpdateMemberAccessRef.ID,
		fieldsToUpdateMemberAccessRef,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}

	return updatedMemberAccessRef, nil
}
