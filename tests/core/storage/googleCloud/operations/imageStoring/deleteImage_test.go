package googlecloudstoragecoreoperationtests

import (
	"context"
	"errors"
	"fmt"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) TestDeleteImageOK() {
	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, fmt.Sprintf("images%s", gcsOperationTestSuite.objectPath),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	gcsOperationTestSuite.mockedGCSObjectHandle.On(
		"Delete", mock.AnythingOfType("*context.timerCtx"),
	).Return(nil).
		Once()

	output, err := gcsOperationTestSuite.basicImageStoringUnderTest.DeleteImage(
		context.Background(),
		fmt.Sprintf("https://storage.googleapis.com/%s/%s",
			gcsOperationTestSuite.bucketName,
			gcsOperationTestSuite.objectPath,
		),
	)

	gcsOperationTestSuite.mockedGCSClient.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSObjectHandle.AssertExpectations(gcsOperationTestSuite.T())

	assert.Equal(
		gcsOperationTestSuite.T(),
		true,
		output,
	)
	assert.Nil(
		gcsOperationTestSuite.T(),
		err,
	)
}

func (gcsOperationTestSuite *GoogleCloudStorageOperationTestSuite) TestDeleteImageDeleteError() {
	gcsOperationTestSuite.mockedGCSClient.
		On(
			"GetObjectHandle", gcsOperationTestSuite.bucketName, fmt.Sprintf("images%s", gcsOperationTestSuite.objectPath),
		).
		Return(gcsOperationTestSuite.mockedGCSObjectHandle, nil).
		Once()

	gcsOperationTestSuite.mockedGCSObjectHandle.On(
		"Delete", mock.AnythingOfType("*context.timerCtx"),
	).Return(errors.New("Some Upstream Error")).
		Once()

	output, err := gcsOperationTestSuite.basicImageStoringUnderTest.DeleteImage(
		context.Background(),
		fmt.Sprintf("https://storage.googleapis.com/%s/%s",
			gcsOperationTestSuite.bucketName,
			gcsOperationTestSuite.objectPath,
		),
	)

	gcsOperationTestSuite.mockedGCSClient.AssertExpectations(gcsOperationTestSuite.T())
	gcsOperationTestSuite.mockedGCSObjectHandle.AssertExpectations(gcsOperationTestSuite.T())

	assert.Equal(
		gcsOperationTestSuite.T(),
		false,
		output,
	)
	assert.Equal(
		gcsOperationTestSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DeleteImageFailed,
			"GoogleCloudStorageBasicOperation",
			errors.New("Some Upstream Error"),
		),
		err,
	)
}
