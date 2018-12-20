package main

import (
	"github.com/xeenhl/spendlog/backend/authService/jwt"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/xeenhl/spendlog/backend/authService/config"
	"github.com/xeenhl/spendlog/backend/authService/handlers"
	"github.com/xeenhl/spendlog/backend/authService/server"
	"github.com/xeenhl/spendlog/backend/authService/storage"
)

func main() {

	c, err := config.LoadConfiguration("./resources/config.json")

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Handle("/sso", negroni.New(negroni.HandlerFunc(handlers.Sso))).Methods("POST")
	router.Handle("/login", negroni.New(negroni.HandlerFunc(handlers.Login))).Methods("POST")
	router.Handle("/logout", negroni.New(negroni.HandlerFunc(handlers.Logout))).Methods("POST")
	router.Handle("/signin", negroni.New(negroni.HandlerFunc(handlers.Signin))).Methods("POST")

	storage := storage.NewMemoryStore()

	s := server.Server{
		Config:       *c,
		Tokenizer:    jwt.NewTokenizer(jwt.NewFileKeyLoader(*c)),
		Router:       router,
		UserStore:    storage,
		SessionStore: storage,
	}

	s.Run()
}
