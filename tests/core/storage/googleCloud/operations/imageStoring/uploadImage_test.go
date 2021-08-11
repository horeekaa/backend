package googlecloudstoragecoreoperationtests

import (
	"context"
	"errors"
	"fmt"
	"io"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlecloudstoragecoretypes "github.com/horeekaa/backend/core/storages/googleCloud/types"
	"github.com/horeekaa/backend/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) TestUploadImageOK() {
	mockedFile := &struct {
		io.Reader
	}{}

	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, mock.AnythingOfType("string"),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	gcsOperationTestSuite.mockedGCSClient.On(
		"CopyWrite", mock.Anything, mock.Anything,
	).Return(int64(1), nil).
		Once()

	gcsOperationTestSuite.mockedGCSObjectHandle.On(
		"NewWriter", mock.AnythingOfType("*context.timerCtx"),
	).Return(gcsOperationTestSuite.mockedGCSWriter).
		Once()

	gcsOperationTestSuite.mockedGCSWriter.On(
		"Close",
	).Return(nil).
		Once()

	output, err := gcsOperationTestSuite.basicImageStoringUnderTest.UploadImage(
		context.Background(),
		model.DescriptivePhotoCategoryOrganizationProfile,
		googlecloudstoragecoretypes.GCSFileUpload{
			File: mockedFile,
		},
	)

	gcsOperationTestSuite.mockedGCSClient.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSObjectHandle.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSWriter.AssertExpectations(gcsOperationTestSuite.T())

	assert.Contains(
		gcsOperationTestSuite.T(),
		output,
		fmt.Sprintf("https://storage.googleapis.com/%s", gcsOperationTestSuite.bucketName),
	)
	assert.Nil(
		gcsOperationTestSuite.T(),
		err,
	)
}

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) TestUploadImageCopyWriteError() {
	mockedFile := &struct {
		io.Reader
	}{}

	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, mock.AnythingOfType("string"),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	gcsOperationTestSuite.mockedGCSClient.On(
		"CopyWrite", mock.Anything, mock.Anything,
	).Return(int64(0), errors.New("Some Upstream Error")).
		Once()

	gcsOperationTestSuite.mockedGCSObjectHandle.On(
		"NewWriter", mock.AnythingOfType("*context.timerCtx"),
	).Return(gcsOperationTestSuite.mockedGCSWriter).
		Once()

	output, err := gcsOperationTestSuite.basicImageStoringUnderTest.UploadImage(
		context.Background(),
		model.DescriptivePhotoCategoryOrganizationProfile,
		googlecloudstoragecoretypes.GCSFileUpload{
			File: mockedFile,
		},
	)

	gcsOperationTestSuite.mockedGCSClient.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSObjectHandle.AssertExpectations(gcsOperationTestSuite.T())

	assert.Zero(
		gcsOperationTestSuite.T(),
		output,
	)
	assert.Equal(
		gcsOperationTestSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.StoringImageFailed,
			"/gcsBasicOperation/uploadImage",
			errors.New("Some Upstream Error"),
		),
		err,
	)
}

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) TestUploadImageWriterCloseError() {
	mockedFile := &struct {
		io.Reader
	}{}

	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, mock.AnythingOfType("string"),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	gcsOperationTestSuite.mockedGCSClient.On(
		"CopyWrite", mock.Anything, mock.Anything,
	).Return(int64(1), nil).
		Once()

	gcsOperationTestSuite.mockedGCSObjectHandle.On(
		"NewWriter", mock.AnythingOfType("*context.timerCtx"),
	).Return(gcsOperationTestSuite.mockedGCSWriter).
		Once()

	gcsOperationTestSuite.mockedGCSWriter.On(
		"Close",
	).Return(errors.New("Some Upstream Error")).
		Once()

	output, err := gcsOperationTestSuite.basicImageStoringUnderTest.UploadImage(
		context.Background(),
		model.DescriptivePhotoCategoryOrganizationProfile,
		googlecloudstoragecoretypes.GCSFileUpload{
			File: mockedFile,
		},
	)

	gcsOperationTestSuite.mockedGCSClient.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSObjectHandle.AssertExpectations(gcsOperationTestSuite.T())

	assert.Zero(
		gcsOperationTestSuite.T(),
		output,
	)
	assert.Equal(
		gcsOperationTestSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClosingImageStoringWriterFailed,
			"/gcsBasicOperation/uploadImage",
			errors.New("Some Upstream Error"),
		),
		err,
	)
}
