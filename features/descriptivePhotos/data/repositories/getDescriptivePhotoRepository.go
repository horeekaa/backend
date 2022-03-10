package descriptivephotodomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getDescriptivePhotoRepository struct {
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	pathIdentity               string
}

func NewGetDescriptivePhotoRepository(
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
) (descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository, error) {
	return &getDescriptivePhotoRepository{
		descriptivePhotoDataSource,
		"GetDescriptivePhotoRepository",
	}, nil
}

func (getDescriptivePhotoRepo *getDescriptivePhotoRepository) Execute(filterFields *model.DescriptivePhotoFilterFields) (*model.DescriptivePhoto, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	descriptivePhoto, err := getDescriptivePhotoRepo.descriptivePhotoDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getDescriptivePhotoRepo.pathIdentity,
			err,
		)
	}

	return descriptivePhoto, nil
}
