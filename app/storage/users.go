package storage

type User struct {
	ID int64
}

type UserRepository interface {
	Find(userID int64) (*User, error)
	CreateUser(user User) error
}

type UserMemoryStorage struct {
	data map[int64]*User
}

func NewUserMemoryStorage() *UserMemoryStorage {
	return &UserMemoryStorage{data: make(map[int64]*User)}
}

func (u *UserMemoryStorage) Find(userID int64) (*User, error) {
	return u.data[userID], nil
}

func (u *UserMemoryStorage) CreateUser(user User) error {
	u.data[user.ID] = &user
	return nil
}
