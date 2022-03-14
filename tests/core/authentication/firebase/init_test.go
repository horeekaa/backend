package firebaseauthcoreclienttests

import (
	"errors"

	firebaseauthcoreclients "github.com/horeekaa/backend/core/authentication/firebase"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) TestInitializeReturnOK() {
	fbAuthSuite.mockedFBApp.
		On("Auth", mock.AnythingOfType("*context.emptyCtx")).
		Return(fbAuthSuite.mockedAuthClient, nil).
		Once()

	fbAuthSuite.mockedFBServerlessClient.On("GetApp").
		Return(fbAuthSuite.mockedFBApp, nil).
		Once()

	fbAuthClient, err := firebaseauthcoreclients.NewFirebaseAuthClient()

	statusEcho, err := fbAuthClient.InitializeClient(fbAuthSuite.mockedFBServerlessClient)
	assert.Nil(fbAuthSuite.T(), err)
	assert.EqualValues(fbAuthSuite.T(), true, statusEcho)

	svlClient, err := fbAuthClient.GetServerlessClient()
	assert.Nil(fbAuthSuite.T(), err)
	assert.Equal(fbAuthSuite.T(), fbAuthSuite.mockedFBServerlessClient, svlClient)

	authClient, err := fbAuthClient.GetAuthClient()
	assert.Nil(fbAuthSuite.T(), err)
	assert.Equal(fbAuthSuite.T(), fbAuthSuite.mockedAuthClient, authClient)

	fbAuthSuite.mockedFBApp.AssertExpectations(fbAuthSuite.T())
	fbAuthSuite.mockedFBServerlessClient.AssertExpectations(fbAuthSuite.T())
}

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) TestInitializeReturnError() {
	fbAuthSuite.mockedFBApp.
		On("Auth", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	fbAuthSuite.mockedFBServerlessClient.On("GetApp").
		Return(fbAuthSuite.mockedFBApp, nil).
		Once()

	fbAuthClient, err := firebaseauthcoreclients.NewFirebaseAuthClient()

	statusEcho, err := fbAuthClient.InitializeClient(fbAuthSuite.mockedFBServerlessClient)
	assert.EqualValues(fbAuthSuite.T(), false, statusEcho)
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"FirebaseAuthClient",
			errors.New("Some Upstream Error"),
		),
		err,
	)

	svlClient, err := fbAuthClient.GetServerlessClient()
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"FirebaseAuthClient",
			nil,
		),
		err,
	)
	assert.Nil(fbAuthSuite.T(), svlClient)

	authClient, err := fbAuthClient.GetAuthClient()
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"FirebaseAuthClient",
			nil,
		),
		err,
	)
	assert.Nil(fbAuthSuite.T(), authClient)

	fbAuthSuite.mockedFBApp.AssertExpectations(fbAuthSuite.T())
	fbAuthSuite.mockedFBServerlessClient.AssertExpectations(fbAuthSuite.T())
}
