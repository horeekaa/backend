package firebaseauthcoreclienttests

import (
	"testing"

	firebaseauthcoremocks "github.com/horeekaa/backend/tests/mocks/core/authentication/firebase/interfaces/wrappers"
	firebaseserverlesscoreclientmocks "github.com/horeekaa/backend/tests/mocks/core/serverless/firebase/interfaces"
	firebaseserverlesscoreappmocks "github.com/horeekaa/backend/tests/mocks/core/serverless/firebase/interfaces/wrappers"
	"github.com/stretchr/testify/suite"
)

type FirebaseAuthCoreClientTestSuite struct {
	suite.Suite
	mockedFBServerlessClient *firebaseserverlesscoreclientmocks.FirebaseServerlessClient
	mockedFBApp              *firebaseserverlesscoreappmocks.FirebaseApp
	mockedAuthClient         *firebaseauthcoremocks.FirebaseAuthClient
}

func (fbAuthSuite *FirebaseAuthCoreClientTestSuite) SetupSuite() {
	fbAuthSuite.mockedFBServerlessClient = &firebaseserverlesscoreclientmocks.FirebaseServerlessClient{}
	fbAuthSuite.mockedFBApp = &firebaseserverlesscoreappmocks.FirebaseApp{}
	fbAuthSuite.mockedAuthClient = &firebaseauthcoremocks.FirebaseAuthClient{}
}

func TestFirebaseAuthCoreClientSuite(t *testing.T) {
	suite.Run(t, new(FirebaseAuthCoreClientTestSuite))
}
