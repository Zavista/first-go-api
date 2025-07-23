package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// APIServer is a simple HTTP server that listens for incoming requests
type APIServer struct {
	listenAddr string
	store      AccountStore
}

// NewAPIServer creates a new APIServer instance with the specified listen address.
// FACTORY pattern
func NewAPIServer(listenAddr string, store AccountStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Start() {
	router := http.NewServeMux()

	router.HandleFunc("/account/", makeHTTPHandleFunc(s.handleAccountRouter))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccountRouter))

	fmt.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

// handleAccountRouter manually creates a router since we want to try without using chi/gin
func (s *APIServer) handleAccountRouter(w http.ResponseWriter, req *http.Request) error {
	path := strings.TrimPrefix(req.URL.Path, "/account") // removes the "/account" from the path
	path = strings.Trim(path, "/")                       // removes leading/trailing slashes

	segments := strings.Split(path, "/") // splits into different segments (ex. /account/1/balance => ["1", "balance"]

	switch len(segments) {
	case 0:
		// /account (base path)
		if req.Method == "POST" {
			return s.handleCreateAccount(w, req)
		}
		return fmt.Errorf("method %s not allowed on /account", req.Method)

	case 1:
		// /account/{id}
		id, err := strconv.Atoi(segments[0])
		if err != nil {
			return fmt.Errorf("invalid account ID: %v", err)
		}

		switch req.Method {
		case "GET":
			return s.handleGetAccount(w, req, id)
		case "PUT":
			return s.handleUpdateAccount(w, req, id)
		case "DELETE":
			return s.handleDeleteAccount(w, req, id)
		default:
			return fmt.Errorf("method %s not allowed on /account/{id}", req.Method)
		}

	case 2:
		// /account/{id}/{action} like /account/1/balance
		id, err := strconv.Atoi(segments[0])
		if err != nil {
			return fmt.Errorf("invalid account ID: %v", err)
		}

		action := segments[1]
		switch action {
		case "balance":
			if req.Method == "GET" {
				return s.handleGetBalance(w, req, id)
			}
		}
	}

	return fmt.Errorf("not found")
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, req *http.Request, id int) error {

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	var createReq CreateAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		log.Printf("failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body")
	}

	created, err := s.store.CreateAccount(&createReq)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, created)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request, id int) error {
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, req *http.Request, id int) error {
	var updateReq UpdateAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		log.Printf("failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body")
	}

	updated, err := s.store.UpdateAccount(id, &updateReq)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updated)
}

func (s *APIServer) handleGetBalance(w http.ResponseWriter, req *http.Request, id int) error {
	balance, err := s.store.GetAccountBalanceByID(id)
	if err != nil {
		return err
	}

	resp := BalanceResponse{
		ID:      id,
		Balance: balance,
	}
	return WriteJSON(w, http.StatusOK, resp)
}

// WriteJSON is a helper function that writes a JSON response with the given status code and data.
// It sets the Content-Type to "application/json" and uses json.Encoder to write the response body.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// apiFunc is a custom function signature that wraps HTTP handlers but returns an error.
// This allows us to centralize error handling using middleware logic
type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

// makeHTTPHandleFunc takes an apiFunc and returns a standard http.HandlerFunc.
// this is necessary since standard http.HandlerFunc does not accept Error in the function signature but we want to handle error outside of the function
// so we handle it here, in one centralized handler location
// btw this is the DECORATOR pattern
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := f(w, req); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
