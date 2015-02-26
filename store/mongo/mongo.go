package mongo

import (
	"github.com/plimble/crud/mongo"
	"github.com/plimble/repass"
	"gopkg.in/mgo.v2"
)

type MongoStore struct {
	*mongo.CRUD
	session *mgo.Session
	db      string
	c       string
}

func NewMongoStore(session *mgo.Session, db string) *MongoStore {
	return &MongoStore{mongo.New(session, db, "repass"), session, db, "repass"}
}

func (m *MongoStore) Get(tokenID string) (*repass.Token, error) {
	session := m.session.Clone()
	defer session.Close()

	var token *repass.Token
	if err := session.DB(m.db).C(m.c).FindId(tokenID).One(&token); err != nil {
		return nil, err
	}

	return token, nil
}
