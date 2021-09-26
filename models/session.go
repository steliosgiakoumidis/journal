package models

import (
	"time"
)

type Session struct {
	Id        int       `json:"id"`
	SubjectId int       `json:"subjectid"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	Tags      []string  `json:"tags"`
	Progress  int       `json:"progress"`
	Date      time.Time `json:"date"`
	Price     int       `json:"price"`

	//TODO add createdAt updatedAt
}

func (s *Session) IsValid() bool {
	if s.Title == "" {
		return false
	}
	if len(s.Tags) == 0 {
		return false
	}
	if s.SubjectId == 0 {
		return false
	}
	if s.Notes == "" {
		return false
	}
	if s.Progress == 0 {
		return false
	}
	if s.Price == 0 {
		return false
	}

	return true
}
