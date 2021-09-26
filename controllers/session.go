package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"journal/models"
	"log"
	"net/http"
	"strconv"
)

type SessionRepository interface {
	GetSessions(id int) ([]models.Session, error)
	GetSession(sessionId int) (*models.Session, error)
	UpdateSession(session models.Session) error
	InsertSession(session models.Session) error
	DeleteSession(sessionId int) error
}

type SessionHandler struct {
	repository SessionRepository
}

func NewSessionHandler(sessionRepository SessionRepository) *SessionHandler {
	return &SessionHandler{
		repository: sessionRepository,
	}
}

func (s *SessionHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	var err error
	var sessions []models.Session
	if sessions, err = s.repository.GetSessions(0); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	if len(sessions) == 0 {
		http.Error(w, http.StatusText(404), 404)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(sessions)
}

func (s *SessionHandler) GetSessionsForSubject(w http.ResponseWriter, r *http.Request) {
	var err error
	var id string
	var sessions []models.Session
	id = extractIdFromUrl(id, r, true)
	subjectId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Validation failed for session id. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if sessions, err = s.repository.GetSessions(subjectId); err != nil {
		log.Println("Get sessions from db failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(sessions)
}
func (s *SessionHandler) GetSessionById(w http.ResponseWriter, r *http.Request) {
	var err error
	var id string
	var session *models.Session
	id = extractIdFromUrl(id, r, false)
	sessionId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Validation failed for session id. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if session, err = s.repository.GetSession(sessionId); err != nil {
		log.Println("Get sessions from db failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(session)
}

func extractIdFromUrl(id string, r *http.Request, forSubjectId bool) string {
	if id = chi.URLParam(r, "sessionid"); forSubjectId {
		id = chi.URLParam(r, "subjectid")
	}
	return id
}

func (s *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {

	var err error
	session := models.Session{}
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil || session.IsValid() == false {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = s.repository.UpdateSession(session); err != nil {
		println("Update session failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func (s *SessionHandler) InsertSession(w http.ResponseWriter, r *http.Request) {
	var err error
	session := models.Session{}
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil || session.IsValid() == false {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = s.repository.InsertSession(session); err != nil {
		println("Insert session failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func (s *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	var err error
	id := chi.URLParam(r, "sessionid")
	sessionId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Validation failed for session id. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = s.repository.DeleteSession(sessionId); err != nil {
		log.Println("Get sessions from db failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}
