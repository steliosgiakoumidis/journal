package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"journal/models"
	"log"
	"net/http"
	"strconv"
)

type TagRepository interface {
	GetAllTags() ([]models.Tag, error)
	InsertTag(tag models.Tag) error
	DeleteTag(idInt int) error
}

type TagHandler struct {
	tagRepository TagRepository
}

func NewTagHandler(repository TagRepository) *TagHandler {
	return &TagHandler{
		tagRepository: repository,
	}
}

func (t *TagHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	var err error
	var tags []models.Tag
	if tags, err = t.tagRepository.GetAllTags(); err != nil {
		log.Fatalln("Database failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (t *TagHandler) InsertTag(w http.ResponseWriter, r *http.Request) {
	var err error
	var tag models.Tag
	tagName := chi.URLParam(r, "name")
	if tagName == "" {
		println(err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = t.tagRepository.InsertTag(tag); err != nil {
		log.Fatalln("Database failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func (t *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var err error
	var idInt int
	id := chi.URLParam(r, "id")
	idInt, err = strconv.Atoi(id)
	if id == "" || err != nil {
		log.Println("Id is missing or malformed. Error: " + err.Error())
		http.Error(w, http.StatusText(400), 400)
	}

	if err = t.tagRepository.DeleteTag(idInt); err != nil {
		log.Println("Delete db operation failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}
