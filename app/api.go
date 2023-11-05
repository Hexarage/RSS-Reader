package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

func NewAPIServer(listenAddress string, strg Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddress,
		storage:    strg,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/RSSAggregator", makeHTTPHandlerFunc(s.handleRequest)) // figure out the json stuff
	router.HandleFunc("/RSSAggregator/{id}", makeHTTPHandlerFunc(s.handleGetRequest))

	log.Println("JSON API Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
func (s *APIServer) handleRequest(w http.ResponseWriter, rq *http.Request) error {

	switch rq.Method {
	case "GET":
		return s.handleGetRequest(w, rq)
	case "POST":
		return s.handleAddRequest(w, rq)
	case "PATCH":

	case "DELETE":
		return s.handleDeleteRequest(w, rq)
	}

	return nil
}

func (s *APIServer) handleGetRequest(w http.ResponseWriter, rq *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteRequest(w http.ResponseWriter, rq *http.Request) error {
	return nil
}

func (s *APIServer) handleAddRequest(w http.ResponseWriter, rq *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
