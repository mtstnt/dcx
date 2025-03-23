package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type File struct {
	Encoding string `json:"Encoding"`
	Content  string `json:"Content"`
}

type Test struct {
	Input  string `json:"Input"`
	Output string `json:"Output"`
}

type SessionRequest struct {
	Files    map[string]File `json:"Files"`
	Language string          `json:"Language"`
	Tests    []Test          `json:"Tests"`
}

type TestResult struct {
	Input          string `json:"Input"`
	ActualOutput   string `json:"ActualOutput"`
	ExpectedOutput string `json:"ExpectedOutput"`
	Status         string `json:"Status"`
}

type SessionStatus struct {
	Status  string       `json:"Status"`
	Results []TestResult `json:"Results,omitempty"`
}

type Session struct {
	ID      string
	Request SessionRequest
	Status  SessionStatus
}

// Sessions map to store session data in memory
var sessions sync.Map

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()
	session := &Session{
		ID:      sessionID,
		Request: req,
		Status: SessionStatus{
			Status: "PENDING",
		},
	}

	sessions.Store(sessionID, session)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"SessionID": sessionID,
	})
}

func GetSessionStatus(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")

	session, ok := sessions.Load(sessionID)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	s := session.(*Session)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(s.Status)
}

func CancelSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")
	sessions.Delete(sessionID)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiRouter := chi.NewRouter()
	apiRouter.Post("/session", CreateSession)
	apiRouter.Delete("/session/{session_id}", CancelSession)
	apiRouter.Get("/session/{session_id}/status", GetSessionStatus)

	r.Mount("/api", apiRouter)
	http.ListenAndServe(":8080", r)
}
