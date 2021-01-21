package dao

import (
	"log"

	"github.com/nikhilpatne/provider-management/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProviderDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "providers"
)

// Establish a connection to database
func (p *ProviderDAO) Connect() {
	session, err := mgo.Dial(p.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(p.Database)
}

// ---------------------------------------------------------------------------------------------------

// Find list of providers
func (p *ProviderDAO) FindAll() ([]models.Provider, error) {
	var providers []models.Provider
	err := db.C(COLLECTION).Find(bson.M{}).All(&providers)
	return providers, err
}

// Find a book by its id
func (p *ProviderDAO) FindById(id string) (models.Provider, error) {
	var provider models.Provider
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&provider)
	return provider, err
}

// Insert a movie into database
func (p *ProviderDAO) Insert(provider models.Provider) error {
	err := db.C(COLLECTION).Insert(&provider)
	return err
}

// Delete an existing book
func (p ProviderDAO) Delete(provider models.Provider) error {
	err := db.C(COLLECTION).Remove(&provider)
	return err
}

// Update an existing book
func (p *ProviderDAO) Update(provider models.Provider) error {
	err := db.C(COLLECTION).UpdateId(provider.ID, &provider)
	return err
}

