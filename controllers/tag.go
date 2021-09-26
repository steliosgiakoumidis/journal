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
	InsertTag(tag string) error
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
	tagName := chi.URLParam(r, "name")
	log.Println("Url path: " + r.URL.Path)
	if tagName == "" {
		log.Println("tag name parameter is missing")
		http.Error(w, http.StatusText(400), 400)
	}

	if err = t.tagRepository.InsertTag(tagName); err != nil {
		log.Println("Database failed")
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(200)
}

func (t *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var err error
	var idInt int
	id := chi.URLParam(r, "tagId")
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
