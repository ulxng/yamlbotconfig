package storage

import "ulxng/yamlbotconf/model"

type UserMemoryStorage struct {
	data map[int64]*model.User
}

func NewUserMemoryStorage() *UserMemoryStorage {
	return &UserMemoryStorage{data: make(map[int64]*model.User)}
}

func (u *UserMemoryStorage) Find(userID int64) (*model.User, error) {
	return u.data[userID], nil
}

func (u *UserMemoryStorage) CreateUser(user model.User) error {
	u.data[user.ID] = &user
	return nil
}
