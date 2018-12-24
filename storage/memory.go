package storage

import (
	"log"
	"os"

	"github.com/rs/xid"

	"github.com/pkg/errors"
	"authService/model"
)

//Memory storage should be used for development and test purpose only
type MemoryStorage struct {
	Users    map[string]model.User
	Sessions map[string]struct{}
}

var logger = log.New(os.Stdout, "[store] ", log.LstdFlags)

func (f *MemoryStorage) Store(u model.User) error {

	logger.Println(f.Users)

	if _, ok := f.Users[u.Credentials.Email]; !ok {
		f.Users[u.Credentials.Email] = u
		return nil
	}

	return errors.New("Usetr Allready exist")
}

func (f *MemoryStorage) GetUserByLogin(login string) (*model.User, error) {

	logger.Println(f.Users)

	if u, ok := f.Users[login]; ok {
		return &u, nil
	}

	return nil, errors.New("No user found")
}

func (f *MemoryStorage) GetUserById(id string) (*model.User, error) {

	logger.Println(f.Users)

	idx, err := xid.FromString(id)

	if err != nil {
		return nil, err
	}

	for _, u := range f.Users {
		if u.Id == idx {
			return &u, nil
		}
	}

	return nil, errors.New("No user found")
}

func (f *MemoryStorage) IsUserExistByLogin(login string) bool {
	_, ok := f.Users[login]
	return ok
}

//SessionStore Implementation

func (f *MemoryStorage) StoreToken(t string) error {
	if _, exist := f.Sessions[t]; !exist {
		f.Sessions[t] = struct{}{}
		return nil
	}
	return errors.New("Sugested Token allready stored in session")
}

func (f *MemoryStorage) UpdateToken(old, new string) error {
	if _, exist := f.Sessions[old]; exist {
		delete(f.Sessions, old)
		f.Sessions[new] = struct{}{}
		return nil
	}
	return f.StoreToken(new)
}

func (f *MemoryStorage) DeleteToken(t string) error {
	if _, exist := f.Sessions[t]; exist {
		delete(f.Sessions, t)
		return nil
	}
	return errors.New("Token has not been found for deleting")
}

func (f *MemoryStorage) IsTokenPresent(t string) bool {
	_, ok := f.Sessions[t]
	return ok
}

func NewMemoryStore() Store {
	return &MemoryStorage{
		Users:    make(map[string]model.User),
		Sessions: make(map[string]struct{}),
	}
}
