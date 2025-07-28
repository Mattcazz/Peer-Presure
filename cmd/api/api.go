package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Mattcazz/Peer-Presure.git/service/comment"
	"github.com/Mattcazz/Peer-Presure.git/service/post"
	"github.com/Mattcazz/Peer-Presure.git/service/user"
	"github.com/Mattcazz/Peer-Presure.git/web"

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

	webSubRouter := router.PathPrefix("/").Subrouter()
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(webSubRouter)

	postStore := post.NewStore(s.db)
	commentStore := comment.NewStore(s.db)
	postHandler := post.NewHandler(postStore, commentStore, userStore)
	postHandler.RegisterRoutes(webSubRouter)

	web.LoadTemplates()
	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
