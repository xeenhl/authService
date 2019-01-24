package model

import (
	"encoding/json"

	"github.com/rs/xid"
)

type User struct {
	Id          xid.ID       `json:"id,omitempty"`
	Credentials *Credentilas `json:"credentials,omitempty"`
	Active      bool         `json:"active"`
	Banned      bool         `json:"banned"`
	SSOData
}

func NewUser(c Credentilas) User {
	return User{
		Id:          xid.New(),
		Credentials: &c,
		Banned:      false,
		Active:      true,
	}
}

type Credentilas struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SSOData struct {
	Valid  *bool   `json:"token_valid,omitempty"`
	Claims []Claim `json:"claims,omitempty"`
}

type Claim struct {
	Key   string
	Value interface{}
}

type AuthError struct {
	ErrorCode string `json:"error"`
	Reason    string `json:"reason"`
}

func (err AuthError) Error() string {
	return string(err.ToBytes())
}

func (err AuthError) ToBytes() []byte {
	b, error := json.Marshal(err)
	if error != nil {
		return []byte("{}")
	}
	return b
}
