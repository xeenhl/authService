package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/xeenhl/authService/config"
	"github.com/xeenhl/authService/jwt"
	"github.com/xeenhl/authService/storage"
)

type Server struct {
	Config       config.Configuration
	Tokenizer    jwt.Tokenazer
	Router       *mux.Router
	UserStore    storage.UserStore
	SessionStore storage.SessionStorage
}

var RunningServer *Server = nil

func (s Server) Run() {

	port := ":" + strconv.Itoa(s.Config.Port)

	n := negroni.Classic()
	n.UseHandler(s.Router)

	RunningServer = &s

	http.ListenAndServe(port, n)
}
