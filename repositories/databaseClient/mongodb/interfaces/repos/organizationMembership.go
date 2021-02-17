package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type OrganizationMembershipRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error)
	Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.OrganizationMembership, error)
	Create(input *model.CreateOrganizationMembership, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error)
	Update(ID interface{}, updateData *model.UpdateOrganizationMembership, operationOptions *mongooperationmodels.OperationOptions) (*model.OrganizationMembership, error)
}
