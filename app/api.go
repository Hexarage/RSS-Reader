package main

import (
	RSSReader "RSS-Reader/internal"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

type AddLinkRequest struct {
	Link string `json:"link"`
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddr: listenAddress,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/RSSAggregator", makeHTTPHandlerFunc(s.handleRequest)) // figure out the json stuff
	router.HandleFunc("/RSSAggregator/Feeds", makeHTTPHandlerFunc(s.handleGetParsedRequest))

	log.Println("JSON API Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleRequest(w http.ResponseWriter, rq *http.Request) error {

	switch rq.Method {
	case "GET":
		return nil
	case "POST":
		return s.handleAddRequest(w, rq)
	case "PATCH":

	case "DELETE":
		return s.handleDeleteRequest(w, rq)
	}

	return nil
}

func (s *APIServer) handleGetParsedRequest(w http.ResponseWriter, rq *http.Request) error {
	var links []string

	// get all of the links from the request

	results := RSSReader.Parse(links)

	return WriteJSON(w, http.StatusOK, results)
}

func (s *APIServer) handleDeleteRequest(w http.ResponseWriter, rq *http.Request) error {
	return nil
}

func (s *APIServer) handleAddRequest(w http.ResponseWriter, rq *http.Request) error {
	linkCreateRequest := new(AddLinkRequest)
	if err := json.NewDecoder(rq.Body).Decode(&linkCreateRequest); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, linkCreateRequest) // TODO: Change it so that it doesn't use storage at all
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
