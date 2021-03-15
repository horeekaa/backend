package servicerepodependencies

import (
	databaseinstancereferences "github.com/horeekaa/backend/repositories/databaseClient/instanceReferences/repos"
	mongodbclients "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongorepos "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseservicerepo "github.com/horeekaa/backend/services/database/repos/"
)

func InitializeMemberAccessRefService() (databaseservicerepointerfaces.MemberAccessRefService, error) {
	memberAccessRefRepoMongo, err := mongorepos.NewMemberAccessRefRepoMongo(mongodbclients.DatabaseClient)
	if err != nil {
		return nil, err
	}

	return databaseservicerepo.NewMemberAccessRefService(databaseinstancereferences.MemberAccessRefRepo{
		Instance: &memberAccessRefRepoMongo,
	})
}
