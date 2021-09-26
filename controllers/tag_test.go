package controllers

import (
	"journal/models"
	"net/http/httptest"
	"testing"
)

var err error
var tagResponse models.Tag
var sliceTagResponse []models.Tag

type FakeTagRepository struct {
}

func (f FakeTagRepository) GetAllTags() ([]models.Tag, error) {
	return sliceTagResponse, err
}
func (f FakeTagRepository) InsertTag(tag string) error {
	return err
}
func (f FakeTagRepository) DeleteTag(idInt int) error {
	return err
}

func TestGettingLinksSuccessful(t *testing.T) {
	err = nil
	sliceTagResponse = []models.Tag{models.Tag{Name: "sadsa"}}
	handler := NewTagHandler(FakeTagRepository{})
	req := httptest.NewRequest("GET", "/tags", nil)
	w := httptest.NewRecorder()
	handler.GetTags(w, req)

	if w.Code != 200 {
		t.Error("Response for malformed request (missing params) should be 400")
	}
}

func TestCreateLinkFailsWhenUrlParamIsMissing(t *testing.T) {
	err = nil
	handler := NewTagHandler(FakeTagRepository{})
	req := httptest.NewRequest("POST", "/tags/", nil)
	w := httptest.NewRecorder()
	handler.InsertTag(w, req)

	if w.Code != 400 {
		t.Error("Response for malformed request (missing params) should be 400")
	}
}
