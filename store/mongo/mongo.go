package mongo

import (
	"github.com/plimble/crud/crudmongo"
	"gopkg.in/mgo.v2"
)

type MongoStore struct {
	*crudmongo.CRUD
	session *mgo.Session
	db      string
	c       string
}

func NewMongoStore(session *mgo.Session, db string) *MongoStore {
	return &MongoStore{crudmongo.New(session, db, "repass"), session, db, "repass"}
}

func (m *MongoStore) Get(tokenID string) (*Token, error) {
	session := m.session.Clone()
	defer session.Close()

	var token *Token
	if err := session.DB(m.db).C(m.c).FindId(tokenID).One(&token); err != nil {
		return nil, err
	}

	return token, nil
}
