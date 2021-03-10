package servicerepodependencies

import (
	databaseinstancereferences "github.com/horeekaa/backend/repositories/databaseClient/instanceReferences/repos"
	mongodbclients "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongorepos "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseservicerepo "github.com/horeekaa/backend/services/database/repos/"
)

func InitializePersonService() (databaseservicerepointerfaces.PersonService, error) {
	personRepoMongo, err := mongorepos.NewPersonRepoMongo(mongodbclients.DatabaseClient)
	if err != nil {
		return nil, err
	}

	return databaseservicerepo.NewPersonService(databaseinstancereferences.PersonRepo{
		Instance: &personRepoMongo,
	})
}
