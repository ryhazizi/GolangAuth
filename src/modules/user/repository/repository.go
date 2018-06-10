package repository

import (
	"golangauth/src/modules/user/model"
)

type UserRepository interface {
	Insert(*model.User) error
	FindAll() (model.Users, error)
	FindByEmail(string) (*model.User, error)
}
