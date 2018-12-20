package storage

import (
	"github.com/xeenhl/spendlog/backend/authService/model"
)

type Store interface {
	UserStore
	SessionStorage
}

type UserStore interface {
	Store(u model.User) error
	GetUserByLogin(login string) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	IsUserExistByLogin(login string) bool
}

type SessionStorage interface {
	StoreToken(t string) error
	UpdateToken(old, new string) error
	IsTokenPresent(t string) bool
	DeleteToken(t string) error
}
