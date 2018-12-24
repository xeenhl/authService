package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"authService/model"
	"authService/server"
)

var logger = log.New(os.Stdout, "[authenticaton] ", log.LstdFlags)

func Login(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	c, err := retriveCredentials(r)
	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Cant retrive credentials"}, http.StatusBadRequest)
		return
	}

	err = isCredentialsFull(c)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	uStore := server.RunningServer.UserStore

	u, err := uStore.GetUserByLogin(c.Email)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Error during retrieving user from storage"}, http.StatusBadRequest)
		return
	}

	logger.Printf("Got user %v", u)

	t := server.RunningServer.Tokenizer
	token, err := t.GenerateToken(u)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Error during token generation"}, http.StatusBadRequest)
		return
	}

	tStore := server.RunningServer.SessionStore
	tStore.StoreToken(token)

	rw.Write([]byte(token))

	logger.Printf("Loging done for user %v, with token %v", c, t)
	logger.Println("Hello from Login")
}

func Logout(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := retriveToken(r)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Didn't get token form request body"}, http.StatusBadRequest)
		return
	}

	store := server.RunningServer.SessionStore

	if store.IsTokenPresent(token) {

		err := store.DeleteToken(token)

		if err != nil {
			prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusOK)
		return
	}

	prepareErrorResponse(rw, model.AuthError{"Unauthorized", "User not Logged in"}, http.StatusUnauthorized)

	logger.Println("Logout done")
}

func Signin(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	c, err := retriveCredentials(r)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Cant retrive credentials"}, http.StatusBadRequest)
		return
	}

	err = isCredentialsFull(c)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	uStore := server.RunningServer.UserStore

	if uStore.IsUserExistByLogin(c.Email) {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "User already exist"}, http.StatusBadRequest)
		return
	}

	u := model.NewUser(*c)

	err = uStore.Store(u)
	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	logger.Printf("User %v", u)
	j, err := json.Marshal(u)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	rw.Write(j)

	logger.Println("Signin done")
}

func Sso(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//TODO
	logger.Println("SSO handler started")
	token, err := retriveToken(r)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Didn't get token form request body"}, http.StatusBadRequest)
		return
	}

	if ok := server.RunningServer.SessionStore.IsTokenPresent(token); !ok {
		prepareErrorResponse(rw, model.AuthError{"Unauthorized", "User not logged in"}, http.StatusUnauthorized)
		return
	}

	logger.Printf("Got token %v", token)

	pt, err := server.RunningServer.Tokenizer.ParceAndVerifyToken(token)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	if !pt.Valid {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "Token in invalid"}, http.StatusBadRequest)
		return
	}

	t := new(model.SSOData)
	t.Valid = pt.Valid

	u, claims := getUserForToken(pt)
	t.Claims = claims

	if u == nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", "User not found"}, http.StatusBadRequest)
		return
	}

	t.Active = u.Active
	t.Banned = u.Banned
	t.UserId = u.Id

	resp, err := json.Marshal(t)

	if err != nil {
		prepareErrorResponse(rw, model.AuthError{"Bad Request", err.Error()}, http.StatusBadRequest)
		return
	}

	rw.Write(resp)

}

func retriveCredentials(r *http.Request) (*model.Credentilas, error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	c := &model.Credentilas{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func isCredentialsFull(c *model.Credentilas) error {

	err := []string{}

	if len(c.Email) < 1 {
		err = append(err, "email")
	}
	if len(c.Password) < 1 {
		err = append(err, "password")
	}

	if len(err) != 0 {
		s := "Login fields value are missing: ["
		for _, e := range err {
			s += "'" + e + "', "
		}
		s = strings.TrimSuffix(s, ", ")
		s += "]"
		return errors.New(s)
	}
	return nil
}

func retriveToken(r *http.Request) (string, error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func prepareErrorResponse(rw http.ResponseWriter, err model.AuthError, status int) {
	rw.WriteHeader(status)
	rw.Write(err.ToBytes())
}

func getUserForToken(pt *jwt.Token) (*model.User, []model.Claim) {
	var u *model.User
	cl := []model.Claim{}
	if claims, ok := pt.Claims.(jwt.MapClaims); ok && pt.Valid {
		logger.Println(claims)
		for key, value := range claims {
			logger.Printf("Parcing [%v, %v]", key, value)
			if key == "sub" {
				if id, ok := value.(string); ok {
					u, _ = server.RunningServer.UserStore.GetUserById(string(id))
				}
			}
			cl = append(cl, model.Claim{Key: key, Value: value})
		}
	}
	return u, cl
}
