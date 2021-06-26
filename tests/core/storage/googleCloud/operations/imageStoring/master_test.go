package googlecloudstoragecoreoperationtests

import (
	"testing"

	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	googlecloudstoragecoreoperations "github.com/horeekaa/backend/core/storages/googleCloud/operations"
	googlecloudstoragecoreoperationmocks "github.com/horeekaa/backend/tests/mocks/core/storages/googleCloud/interfaces/init"
	googlecloudstoragecoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/storages/googleCloud/interfaces/wrappers"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GoogleCloudStorageOperationTestSuite struct {
	suite.Suite
	mockedGCSClient       *googlecloudstoragecoreoperationmocks.GoogleCloudStorageClient
	mockedGCSObjectHandle *googlecloudstoragecoreoperationwrappermocks.GCSObjectHandle
	mockedGCSWriter       *googlecloudstoragecoreoperationwrappermocks.GCSWriter
	bucketName            string
	objectPath            string

	basicImageStoringUnderTest googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation
}

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) SetupTest() {
	gcsOperationTestSuite.mockedGCSWriter = &googlecloudstoragecoreoperationwrappermocks.GCSWriter{}
	gcsOperationTestSuite.mockedGCSObjectHandle = &googlecloudstoragecoreoperationwrappermocks.GCSObjectHandle{}
	gcsOperationTestSuite.mockedGCSClient = &googlecloudstoragecoreoperationmocks.GoogleCloudStorageClient{}
	gcsOperationTestSuite.bucketName = "MyBucket"
	gcsOperationTestSuite.objectPath = "/my/object/path"

	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, mock.AnythingOfType("string"),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	basicGCSImageStoring, _ := googlecloudstoragecoreoperations.NewGCSBasicImageStoringOperation(
		gcsOperationTestSuite.mockedGCSClient,
		gcsOperationTestSuite.bucketName,
	)

	gcsOperationTestSuite.basicImageStoringUnderTest = basicGCSImageStoring
}

func TestGCSOperationSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageOperationTestSuite))
}
