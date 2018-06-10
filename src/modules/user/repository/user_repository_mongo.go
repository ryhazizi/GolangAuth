package repository

import (
	"golangauth/src/modules/user/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userRepositoryMongo struct {
	db         *mgo.Database
	collection string
}

func NewUserRepositoryMongo(db *mgo.Database, collection string) *userRepositoryMongo {
	return &userRepositoryMongo{
		db:         db,
		collection: collection,
	}
}

func (r *userRepositoryMongo) Insert(user *model.User) error {
	err := r.db.C(r.collection).Insert(user)
	return err
}

func (r *userRepositoryMongo) FindAll() (model.Users, error) {
	var users model.Users

	err := r.db.C(r.collection).Find(bson.M{}).All(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryMongo) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.C(r.collection).Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
