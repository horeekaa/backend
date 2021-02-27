package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type personService struct {
	personRepo *mongorepointerfaces.PersonRepoMongo
}

func NewPersonService(personRepo mongorepointerfaces.PersonRepoMongo) (databaseservicerepointerfaces.PersonService, error) {
	return &personService{
		&personRepo,
	}, nil
}

func (personSvc *personService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Person, chan error) {
	personChn := make(chan *model.Person)
	errorChn := make(chan error)

	go func() {
		person, err := (*personSvc.personRepo).FindByID(ID, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/personService/FindByID",
				&err,
			)
			return
		}

		personChn <- person
	}()

	return personChn, errorChn
}

func (personSvc *personService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Person, chan error) {
	personChn := make(chan *model.Person)
	errorChn := make(chan error)

	go func() {
		person, err := (*personSvc.personRepo).FindOne(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/personService/FindOne",
				&err,
			)
			return
		}

		personChn <- person
	}()

	return personChn, errorChn
}

func (personSvc *personService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.Person, chan error) {
	personsChn := make(chan []*model.Person)
	errorChn := make(chan error)

	go func() {
		persons, err := (*personSvc.personRepo).Find(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/personService/Find",
				&err,
			)
			return
		}

		personsChn <- persons
	}()

	return personsChn, errorChn
}

func (personSvc *personService) Create(input *model.CreatePerson, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Person, chan error) {
	personChn := make(chan *model.Person)
	errorChn := make(chan error)

	go func() {
		person, err := (*personSvc.personRepo).Create(input, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/personService/Create",
				&err,
			)
			return
		}

		personChn <- person
	}()

	return personChn, errorChn
}

func (personSvc *personService) Update(ID interface{}, updateData *model.UpdatePerson, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Person, chan error) {
	personChn := make(chan *model.Person)
	errorChn := make(chan error)

	go func() {
		person, err := (*personSvc.personRepo).Update(ID, updateData, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/personService/Update",
				&err,
			)
			return
		}

		personChn <- person
	}()

	return personChn, errorChn
}
