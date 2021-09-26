package db

import (
	"database/sql"
	"fmt"
	"journal/models"
	"log"
	"strconv"
)

type SessionRepository struct {
	TagRepository TagRepository
}

func (s SessionRepository) GetSessions(subjectId int) ([]models.Session, error) {
	var rows *sql.Rows
	var err error
	if rows, err = DbConn.Query("SELECT * FROM session order by id desc"); subjectId != 0 {
		rows, err = DbConn.Query("SELECT * FROM session where subject_id=$1", subjectId)
	}
	if err != nil {
		fmt.Println("An error occured when getting all sessions. Error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		session := models.Session{}
		if err = rows.Scan(&session.Id, &session.Title, &session.Notes, &session.Progress, &session.Date, &session.Price, &session.SubjectId); err != nil {
			println("An error occured when assigning subjects. Error: " + err.Error())
			return nil, err
		}
		sessions = append(sessions, session)
	}

	for _, v := range sessions {
		if v.Tags, err = s.getTagsForSession(v.Id); err != nil {
			fmt.Println("An error occured when getting tags for session_id: " + strconv.Itoa(v.Id) + ". Error: " + err.Error())
			return nil, err
		}
	}

	return sessions, nil
}

func (s SessionRepository) GetSession(id int) (*models.Session, error) {
	var err error
	var row *sql.Row
	row = DbConn.QueryRow("SELECT * FROM session where id=$1", id)
	if row.Err() != nil {
		println("An error occured when getting session with id " + strconv.Itoa(id) + ". Error: " + err.Error())
		return nil, err
	}
	var session *models.Session
	err = row.Scan(&session.Id, &session.Title, &session.Notes, &session.Progress, &session.Date, &session.Price, &session.SubjectId)
	switch err {
	case sql.ErrNoRows:
		log.Println("Session with id " + strconv.Itoa(id) + " not found. Error: " + err.Error())
		return nil, err
	case nil:
		log.Println("subject found")
	default:
		log.Println("Unexpected error. Error: " + err.Error())
	}

	if session.Tags, err = s.getTagsForSession(session.Id); err != nil {
		fmt.Println("An error occured when getting tags for session_id: " + strconv.Itoa(id) + ". Error: " + err.Error())
		return nil, err
	}

	return session, nil
}

func (s SessionRepository) UpdateSession(session models.Session) error {
	var err error
	if _, err = DbConn.Exec("UPDATE subject set 'title'=$1,'notes'=$2,'progress'=$3,'date'=$4,'price'=$5,'subject_id'=$6 Where 'id'=$7", session.Title, session.Notes, session.Progress, session.Date, session.Price, session.SubjectId, session.SubjectId); err != nil {
		return err
	}

	if err = s.TagRepository.deleteSessionTags(session.Id); err != nil {
		return err
	}

	if err = s.TagRepository.insertSessionTags(session.Id, session.Tags); err != nil {
		return err
	}

	return nil
}

func (s SessionRepository) InsertSession(session models.Session) error {
	var err error
	if _, err = DbConn.Exec("INSERT INTO subject (title, notes, progress, date, price, subject_id) VALUES ($1,$2, $3,$4,$5,$6)", session.Title, session.Notes, session.Progress, session.Date, session.Progress, session.SubjectId); err != nil {
		return err
	}

	if err = s.TagRepository.insertSessionTags(session.Id, session.Tags); err != nil {
		return err
	}

	return nil
}

func (s SessionRepository) DeleteSession(sessionId int) error {
	var err error
	if err = s.TagRepository.deleteSessionTags(sessionId); err != nil {
		return err
	}

	if _, err = DbConn.Exec("DELETE FROM session where id=$1", sessionId); err != nil {
		return err
	}

	return nil
}

func (s SessionRepository) getTagsForSession(session_id int) ([]string, error) {
	var err error
	var tags []models.Tag
	if tags, err = s.TagRepository.GetSessionTags(session_id); err != nil {
		log.Println("tags not found for session_id: " + strconv.Itoa(session_id) + ". Error: " + err.Error())
		return nil, err
	}

	var tag_list []string
	for _, v := range tags {
		tag_list = append(tag_list, v.Name)
	}

	return tag_list, nil
}
