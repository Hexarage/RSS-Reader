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
	storage    Storage
}

type AddLinkRequest struct {
	Link string `json:"link"`
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
	router.HandleFunc("/RSSAggregator/{id}", makeHTTPHandlerFunc(s.handleGetRequestById))
	router.HandleFunc("/RSSAggregator/Feeds", makeHTTPHandlerFunc(s.handleGetParsedRequest))

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

func (s *APIServer) handleGetRequestById(w http.ResponseWriter, rq *http.Request) error {
	id := mux.Vars(rq)["id"]
	links, err := s.storage.GetFeedById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, links)
}

func (s *APIServer) handleGetParsedRequest(w http.ResponseWriter, rq *http.Request) error {
	links, err := s.storage.GetAllFeeds()
	if err != nil {
		return err
	}

	results := RSSReader.Parse(links)

	return WriteJSON(w, http.StatusOK, results)
}

func (s *APIServer) handleGetRequest(w http.ResponseWriter, rq *http.Request) error {
	links, err := s.storage.GetAllFeeds()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, links)
}

func (s *APIServer) handleDeleteRequest(w http.ResponseWriter, rq *http.Request) error {
	return nil
}

func (s *APIServer) handleAddRequest(w http.ResponseWriter, rq *http.Request) error {
	linkCreateRequest := new(AddLinkRequest)
	if err := json.NewDecoder(rq.Body).Decode(&linkCreateRequest); err != nil {
		return err
	}

	if err := s.storage.AddFeed(linkCreateRequest.Link); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, linkCreateRequest) // TODO: Change storage API so that we get the ID of the link at least
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
