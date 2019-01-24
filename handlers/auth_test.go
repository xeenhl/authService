package handlers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/rs/xid"

	"github.com/pkg/errors"

	"authService/server"

	"authService/model"

	jwt "github.com/dgrijalva/jwt-go"
	"authService/storage"

	"github.com/urfave/negroni"
	"authService/handlers"
)

type TestTokenizer struct {
	secret string
}

const testToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzgyODg3NTksImlhdCI6MTU0MjI4ODc1OSwic3ViIjowfQ.loxAMriJDpfGpMHWTbYo3QeA9ZaE7xwKvt_vpZqU2YY"

var tokenValid = true
var userId, _ = xid.FromString("bfra5o2cc8imh64se1s0")

func (t *TestTokenizer) GenerateToken(u *model.User) (string, error) {

	return testToken, nil

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": u.Id,
	// 	"iat": time.Now().Unix(),
	// 	"exp": time.Now().Add(time.Hour * time.Duration(10000)).Unix(),
	// })

	// s, err := token.SignedString([]byte(t.secret))

	// if err != nil {
	// 	return "", errors.New("No Token has been generated")
	// }

	// return s, nil
}

func (t *TestTokenizer) ParceAndVerifyToken(token string) (*jwt.Token, error) {

	if token == testToken {

		claims := jwt.MapClaims{
			"ExpiresAt": 15000,
			"Issuer":    "test",
			"sub":       userId.String(),
		}

		tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tok.Valid = tokenValid

		return tok, nil
	}

	return nil, errors.New("Wrong test token")

	// return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte(t.secret), nil
	// })

}

var s = storage.NewMemoryStore()

var testSetToDefault = func(s storage.Store) {
	tokenValid = true
	if store, ok := s.(*storage.MemoryStorage); ok {
		store.Users = make(map[string]model.User)
		store.Sessions = make(map[string]struct{})
	}
}

func TestSignin(t *testing.T) {

	tests := []struct {
		description  string
		requestBody  string
		expestedBody string
		expectedCode int
		storeInitter func(s *storage.MemoryStorage)
	}{
		{
			description:  "Should return cant retreave credentials for unprocesable json",
			requestBody:  `"email":"test@gaml.com","password":"qwerty"`,
			expestedBody: `{"error":"Bad Request","reason":"Cant retrive credentials"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for registered user",
			requestBody:  `{"email":"test@gmail.com","password":"qwerty"}`,
			expestedBody: `{"error":"Bad Request","reason":"User already exist"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {
				s.Users["test@gmail.com"] = model.User{}
			},
		},
		{
			description:  "Should return bad request for no email passing",
			requestBody:  `{"email":"","password":"qwerty"}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['email']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for no password passing",
			requestBody:  `{"email":"test@gaml.com","password":""}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['password']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for no email and password passing",
			requestBody:  `{"email":"","password":""}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['email', 'password']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
	}

	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Signin)))
	defer ts.Close()

	for _, tc := range tests {

		t.Run(tc.description, func(t *testing.T) {
			testSetToDefault(s)
			store, ok := s.(*storage.MemoryStorage)
			if ok {
				tc.storeInitter(store)
			}
			res, err := http.Post(ts.URL+"/signin", "application/json", strings.NewReader(tc.requestBody))

			if err != nil {
				t.Error(err.Error())
			}

			if res != nil {
				defer res.Body.Close()
			}

			b, err := ioutil.ReadAll(res.Body)
			str := string(b)

			if ok := IsEqualJson(str, tc.expestedBody); !ok {
				t.Errorf("Wrong response body. Expected: [%v] Actual: [%v]", tc.expestedBody, str)
			}

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Wrong response status code. Expected: [%v] Actual: [%v]", tc.expectedCode, res.StatusCode)
			}
		})
	}
}

func TestSigninCreateNewUser(t *testing.T) {
	test := struct {
		description  string
		requestBody  string
		expestedBody string
		expectedCode int
		storeInitter func(s *storage.MemoryStorage)
	}{
		description:  "Should create new user for notexistion credentials",
		requestBody:  `{"email":"test@gaml.com","password":"qwerty"}`,
		expestedBody: `{"id":"%v","credentials":{"email":"test@gaml.com","password":"qwerty"},"active":true,"banned":false}`,
		expectedCode: 200,
		storeInitter: func(s *storage.MemoryStorage) {

		},
	}

	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Signin)))
	defer ts.Close()

	testSetToDefault(s)
	store, ok := s.(*storage.MemoryStorage)
	if ok {
		test.storeInitter(store)
	}
	res, err := http.Post(ts.URL+"/signin", "application/json", strings.NewReader(test.requestBody))

	if err != nil {
		t.Error(err.Error())
	}

	if res != nil {
		defer res.Body.Close()
	}

	b, err := ioutil.ReadAll(res.Body)
	str := string(b)

	//We asume that user id generated by xid lib is OK so we just take it from response and put to expected result
	responseUser := &model.User{}
	json.Unmarshal(b, responseUser)

	test.expestedBody = fmt.Sprintf(test.expestedBody, responseUser.Id.String())

	if ok := IsEqualJson(str, test.expestedBody); !ok {
		t.Errorf("Wrong response body. Expected: [%v] Actual: [%v]", test.expestedBody, str)
	}

	if res.StatusCode != test.expectedCode {
		t.Errorf("Wrong response status code. Expected: [%v] Actual: [%v]", test.expectedCode, res.StatusCode)
	}
}
func TestSso(t *testing.T) {

	tests := []struct {
		description  string
		requestBody  string
		expestedBody string
		expectedCode int
		testIniter   func(s *storage.MemoryStorage)
	}{
		{
			description:  "Should return Bad request for valid token while there is no user in storage",
			requestBody:  testToken,
			expestedBody: `{"error":"Bad Request","reason":"User not found"}`,
			expectedCode: 400,
			testIniter: func(s *storage.MemoryStorage) {
				s.Sessions[testToken] = struct{}{}
			},
		},
		{
			description:  "Should return Unauthorized for not logged in user",
			requestBody:  testToken,
			expestedBody: `{"error":"Unauthorized","reason":"User not logged in"}`,
			expectedCode: 401,
			testIniter: func(s *storage.MemoryStorage) {
				s.Users["user"] = model.User{Id: userId}
			},
		},
		//This test could be unstable becouse of Claims ordering
		{
			description:  "Should return credentials for provided valid token",
			requestBody:  testToken,
			expestedBody: `{"id":"bfra5o2cc8imh64se1s0","active":false,"banned":false,"token_valid":true,"Claims":[{"Key":"ExpiresAt","Value":15000},{"Key":"Issuer","Value":"test"},{"Key":"sub","Value":"bfra5o2cc8imh64se1s0"}]}`,
			expectedCode: 200,
			testIniter: func(s *storage.MemoryStorage) {
				s.Users["user"] = model.User{
												Id: userId,
												Active: false,
												Banned: false,
											}
				s.Sessions[testToken] = struct{}{}
			},
		},
		{
			description:  "Should return Bad request for invalid token",
			requestBody:  testToken,
			expestedBody: `{"error":"Bad Request","reason":"Token in invalid"}`,
			expectedCode: 400,
			testIniter: func(s *storage.MemoryStorage) {
				s.Users["user"] = model.User{Id: userId}
				s.Sessions[testToken] = struct{}{}
				tokenValid = false
			},
		},
	}

	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s
	server.RunningServer.Tokenizer = &TestTokenizer{secret: "my_test_sercert"}

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Sso)))
	defer ts.Close()

	for _, tc := range tests {

		t.Run(tc.description, func(t *testing.T) {
			testSetToDefault(s)
			store, ok := s.(*storage.MemoryStorage)
			if ok {
				tc.testIniter(store)
			}
			res, err := http.Post(ts.URL+"/sso", "application/json", strings.NewReader(tc.requestBody))

			if err != nil {
				t.Error(err.Error())
			}

			if res != nil {
				defer res.Body.Close()
			}

			b, err := ioutil.ReadAll(res.Body)
			str := string(b)

			if ok := IsEqualJson(str, tc.expestedBody); !ok {
				t.Errorf("Wrong response body. Expected: [%v] Actual: [%v]", tc.expestedBody, str)
			}

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Wrong response status code. Expected: [%v] Actual: [%v]", tc.expectedCode, res.StatusCode)
			}
		})
	}

}

func TestLogin(t *testing.T) {

	tests := []struct {
		description  string
		requestBody  string
		expestedBody string
		expectedCode int
		storeInitter func(s *storage.MemoryStorage)
	}{
		{
			description:  "Should return token for provided credentilas",
			requestBody:  `{"email":"test@gmail.com","password":"qwerty"}`,
			expestedBody: testToken,
			expectedCode: 200,
			storeInitter: func(s *storage.MemoryStorage) {
				s.Users["test@gmail.com"] = model.User{}
			},
		},
		{
			description:  "Should error is there is no user in storage",
			requestBody:  `{"email":"test@gmail.com","password":"qwerty"}`,
			expestedBody: `{"error":"Bad Request","reason":"Error during retrieving user from storage"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return cant retreave credentials for unprocesable json",
			requestBody:  `"email":"test@gaml.com","password":"qwerty"`,
			expestedBody: `{"error":"Bad Request","reason":"Cant retrive credentials"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for no email passing",
			requestBody:  `{"email":"","password":"qwerty"}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['email']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for no password passing",
			requestBody:  `{"email":"test@gaml.com","password":""}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['password']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should return bad request for no email and password passing",
			requestBody:  `{"email":"","password":""}`,
			expestedBody: `{"error":"Bad Request","reason":"Login fields value are missing: ['email', 'password']"}`,
			expectedCode: 400,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
	}

	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s
	server.RunningServer.Tokenizer = &TestTokenizer{secret: "my_test_sercert"}

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Login)))
	defer ts.Close()

	for _, tc := range tests {

		t.Run(tc.description, func(t *testing.T) {
			testSetToDefault(s)
			store, ok := s.(*storage.MemoryStorage)
			if ok {
				tc.storeInitter(store)
			}
			res, err := http.Post(ts.URL+"/login", "application/json", strings.NewReader(tc.requestBody))

			if err != nil {
				t.Error(err.Error())
			}

			if res != nil {
				defer res.Body.Close()
			}

			b, err := ioutil.ReadAll(res.Body)
			str := string(b)

			if ok := IsEqualJson(str, tc.expestedBody); !ok {
				t.Errorf("Wrong response body. Expected: [%v] Actual: [%v]", tc.expestedBody, str)
			}

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Wrong response status code. Expected: [%v] Actual: [%v]", tc.expectedCode, res.StatusCode)
			}
		})
	}

}

func TestLogout(t *testing.T) {

	tests := []struct {
		description  string
		requestBody  string
		expestedBody string
		expectedCode int
		storeInitter func(s *storage.MemoryStorage)
	}{
		{
			description:  "Should Return unathorized for not logged user",
			requestBody:  testToken,
			expestedBody: `{"error":"Unauthorized","reason":"User not Logged in"}`,
			expectedCode: 401,
			storeInitter: func(s *storage.MemoryStorage) {

			},
		},
		{
			description:  "Should Return OK for seccess user logout",
			requestBody:  testToken,
			expestedBody: "",
			expectedCode: 200,
			storeInitter: func(s *storage.MemoryStorage) {
				s.Users["test@gmail.com"] = model.User{Id: userId}
				s.Sessions[testToken] = struct{}{}
			},
		},
	}

	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Logout)))
	defer ts.Close()

	for _, tc := range tests {

		t.Run(tc.description, func(t *testing.T) {
			testSetToDefault(s)
			store, ok := s.(*storage.MemoryStorage)
			if ok {
				tc.storeInitter(store)
			}
			res, err := http.Post(ts.URL+"/signin", "application/json", strings.NewReader(tc.requestBody))

			if err != nil {
				t.Error(err.Error())
			}

			if res != nil {
				defer res.Body.Close()
			}

			b, err := ioutil.ReadAll(res.Body)
			str := string(b)

			if ok := IsEqualJson(str, tc.expestedBody); !ok {
				t.Errorf("Wrong response body. Expected: [%v] Actual: [%v]", tc.expestedBody, str)
			}

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Wrong response status code. Expected: [%v] Actual: [%v]", tc.expectedCode, res.StatusCode)
			}
		})
	}

}

var r *http.Response
var e error

func BenchmarkLogin(b *testing.B) {

	testSetToDefault(s)
	server.RunningServer = &server.Server{}
	server.RunningServer.UserStore = s
	server.RunningServer.SessionStore = s
	body := `{"email":"test@gaml.com","password":"qwerty"}`
	var res *http.Response
	var err error

	ts := httptest.NewServer(negroni.New(negroni.HandlerFunc(handlers.Signin)))
	defer ts.Close()

	for n := 0; n < b.N; n++ {
		res, err = http.Post(ts.URL+"/signin", "application/json", strings.NewReader(body))
	}

	r = res
	e = err
}

func IsEqualJson(s1, s2 string) bool {

	if s1 == s2 {
		return true
	}

	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(s2), &o2)

	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
