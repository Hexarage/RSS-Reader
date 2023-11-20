package main

import (
	RSSReader "RSS-Reader/internal"
	"encoding/json"
	"fmt"
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
	log.Printf("Starting server on address: %v\n", listenAddress)
	return &APIServer{
		listenAddr: listenAddress,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/RSSAggregator", makeHTTPHandlerFunc(s.handleRequest))

	log.Println("JSON API Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleRequest(w http.ResponseWriter, rq *http.Request) error {
	if rq.Method != "GET" {
		return fmt.Errorf("unsupported method")
	}

	var req inputJSON
	if err := json.NewDecoder(rq.Body).Decode(&req); err != nil {
		return err
	}

	if len(req.Links) == 0 {
		return fmt.Errorf("too few elements")
	}
	var result returnJSON
	result.Items = RSSReader.Parse(req.Links)

	if result.Items == nil {
		return WriteJSON(w, http.StatusInternalServerError, nil) // TODO: there's probably more adequate response, change it to that
	}

	return WriteJSON(w, http.StatusOK, result)
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

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc { // This is in case we want to add further functionality
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func slicesAreSame(first []RSSReader.RSSItem, second []RSSReader.RSSItem) bool { // TODO: Make it not care if elements are not in same order
	if len(first) != len(second) {
		return false
	}

	for i, e := range first {
		if e != second[i] {
			return false
		}
	}

	return true
}

type inputJSON struct {
	Links []string `json:"links"`
}

type returnJSON struct {
	Items []RSSReader.RSSItem `json:"items"`
}
