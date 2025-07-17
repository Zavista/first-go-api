package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJson is a helper function that writes a JSON response with the given status code and data.
// It sets the Content-Type to "application/json" and uses json.Encoder to write the response body.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
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
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := f(w, req); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

// APIServer is a simple HTTP server that listens for incoming requests
type APIServer struct {
	listenAddr string
}

// NewAPIServer creates a new APIServer instance with the specified listen address.
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Start() {
	router := http.NewServeMux()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	fmt.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleAccount(w http.ResponseWriter, req *http.Request) error {
	if req.Method == "GET" {
		return s.handleGetAccount(w, req)
	}

	if req.Method == "POST" {
		return s.handleCreateAccount(w, req)
	}

	if req.Method == "DELETE" {
		return s.handleDeleteAccount(w, req)
	}

	return fmt.Errorf("method not allowed %s", req.Method)

}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	return nil

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	return nil

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	return nil

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, req *http.Request) error {
	return nil

}
