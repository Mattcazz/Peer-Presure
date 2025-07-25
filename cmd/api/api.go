package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Mattcazz/Peer-Presure.git/service/user"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
