package storage

import "ulxng/yamlbotconf/model"

type UserRepository interface {
	Find(userID int64) (*model.User, error)
	CreateUser(user model.User) error
}
