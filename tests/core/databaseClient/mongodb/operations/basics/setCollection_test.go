package mongodbcoreoperationtests

import "github.com/stretchr/testify/assert"

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestGetSetCollection() {
	name := mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()

	mongodbBscOpSuite.mockedMongoClient.AssertExpectations(mongodbBscOpSuite.T())
	assert.EqualValues(mongodbBscOpSuite.T(), "MyCollection", name)
}
