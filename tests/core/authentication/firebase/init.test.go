package firebaseauthcoreclienttests

import (
	"errors"

	firebaseauthcoreclients "github.com/horeekaa/backend/core/authentication/firebase"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebaseauthcoreclientfixtures "github.com/horeekaa/backend/tests/fixtures/core/authentication/firebase"
	firebaseserverlesscoreclientmocks "github.com/horeekaa/backend/tests/mocks/core/serverless/firebase/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FirebaseAuthCoreClientTestSuite struct {
	suite.Suite
	mockedFBServerlessClient *firebaseserverlesscoreclientmocks.FirebaseServerlessClient
}

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) SetupSuite() {
	fbAuthSuite.mockedFBServerlessClient = &firebaseserverlesscoreclientmocks.FirebaseServerlessClient{}
}

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) TestInitializeReturnOK() {
	fbAuthSuite.mockedFBServerlessClient.On("GetApp").
		Return(&firebaseauthcoreclientfixtures.FirebaseAppAuthOKFixture{}, nil).
		Once()

	fbAuthClient, err := firebaseauthcoreclients.NewFirebaseAuthClient()

	statusEcho, err := fbAuthClient.InitializeClient(fbAuthSuite.mockedFBServerlessClient)
	assert.Nil(fbAuthSuite.T(), err)
	assert.EqualValues(fbAuthSuite.T(), true, statusEcho)

	client, err := fbAuthClient.GetServerlessClient()
	assert.Nil(fbAuthSuite.T(), err)
	assert.NotNil(fbAuthSuite.T(), client)

	authClient, err := fbAuthClient.GetAuthClient()
	assert.Nil(fbAuthSuite.T(), err)
	assert.NotNil(fbAuthSuite.T(), authClient)

	fbAuthSuite.mockedFBServerlessClient.AssertExpectations(fbAuthSuite.T())
}

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) TestInitializeReturnError() {
	fbAuthSuite.mockedFBServerlessClient.On("GetApp").
		Return(&firebaseauthcoreclientfixtures.FirebaseAppAuthErrorFixture{}, nil).
		Once()

	fbAuthClient, err := firebaseauthcoreclients.NewFirebaseAuthClient()

	statusEcho, err := fbAuthClient.InitializeClient(fbAuthSuite.mockedFBServerlessClient)
	assert.EqualValues(fbAuthSuite.T(), false, statusEcho)
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"/newFirebaseAuthentication",
			errors.New("Some Upstream Error"),
		),
		err,
	)

	client, err := fbAuthClient.GetServerlessClient()
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			nil,
		),
		err,
	)
	assert.Nil(fbAuthSuite.T(), client)

	authClient, err := fbAuthClient.GetAuthClient()
	assert.Equal(fbAuthSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			nil,
		),
		err,
	)
	assert.Nil(fbAuthSuite.T(), authClient)

	fbAuthSuite.mockedFBServerlessClient.AssertExpectations(fbAuthSuite.T())
}
