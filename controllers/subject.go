package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"journal/models"
	"log"
	"net/http"
	"strconv"
)

type SubjectRepository interface {
	GetSubjects() ([]models.Subject, error)
	GetSubject(idInt int) (*models.Subject, error)
	UpdateSubject(subject models.Subject) error
	InsertSubject(subject models.Subject) error
	DeleteSubject(idInt int) error
}

func NewSubjectHandler(repository SubjectRepository) *SubjectHandler {
	return &SubjectHandler{
		subjectRepository: repository,
	}
}

type SubjectHandler struct {
	subjectRepository SubjectRepository
}

func (s *SubjectHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	var err error
	var subjects []models.Subject
	if subjects, err = s.subjectRepository.GetSubjects(); err != nil {
		log.Println("Database failed")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func (s *SubjectHandler) GetSubject(w http.ResponseWriter, r *http.Request) {
	var err error
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var idInt int
	if idInt, err = strconv.Atoi(id); err != nil {
		log.Println("Id cannot be parsed")
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var subject *models.Subject
	if subject, err = s.subjectRepository.GetSubject(idInt); err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

func (s *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	var err error
	subject := models.Subject{}
	err = json.NewDecoder(r.Body).Decode(&subject)
	if err != nil || subject.IsValid() == false {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if err = s.subjectRepository.UpdateSubject(subject); err != nil {
		log.Println("Update failed")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *SubjectHandler) InsertSubject(w http.ResponseWriter, r *http.Request) {
	var err error
	subject := models.Subject{}
	err = json.NewDecoder(r.Body).Decode(&subject)
	if err != nil || subject.IsValid() == false {
		log.Println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = s.subjectRepository.InsertSubject(subject); err != nil {
		log.Println("Insert failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func (s *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	var err error
	var idInt int

	id := chi.URLParam(r, "subjectid")
	idInt, err = strconv.Atoi(id)
	if id == "" || err != nil {
		log.Println("Id is missing or malformed. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = s.subjectRepository.DeleteSubject(idInt); err != nil {
		log.Println("Delete db operation failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}
