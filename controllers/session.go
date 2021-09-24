package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"journal/db"
	"journal/models"
	"log"
	"net/http"
	"strconv"
)

func GetSessions(w http.ResponseWriter, r *http.Request) {
	var err error
	var sessions []models.Session
	if sessions, err = db.GetSessions(); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	if len(sessions) == 0 {
		http.Error(w, http.StatusText(404), 404)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(sessions)
}

func GetSessionForSubject(w http.ResponseWriter, r *http.Request) {
	getSession(w, r, true)
}
func GetSessionForSessionId(w http.ResponseWriter, r *http.Request) {
	getSession(w, r, false)
}

func getSession(w http.ResponseWriter, r *http.Request, forSubjectId bool) {
	var err error
	var id string
	var sessions *models.Session
	id = extractIdFromUrl(id, r, forSubjectId)
	subjectId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Validation failed for session id. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if sessions, err = db.GetSession(subjectId, forSubjectId); err != nil {
		log.Println("Get sessions from db failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(sessions)
}

func extractIdFromUrl(id string, r *http.Request, forSubjectId bool) string {
	if id = chi.URLParam(r, "sessionid"); forSubjectId {
		id = chi.URLParam(r, "subjectid")
	}
	return id
}

func UpdateSession(w http.ResponseWriter, r *http.Request) {

	var err error
	session := models.Session{}
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil || session.IsValid() == false {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = db.UpdateSession(session); err != nil {
		println("Update session failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func InsertSession(w http.ResponseWriter, r *http.Request) {
	var err error
	session := models.Session{}
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil || session.IsValid() == false {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = db.InsertSession(session); err != nil {
		println("Insert session failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	var err error
	id := chi.URLParam(r, "sessionid")
	sessionId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Validation failed for session id. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = db.DeleteSession(sessionId); err != nil {
		log.Println("Get sessions from db failed. Error: " + err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}
