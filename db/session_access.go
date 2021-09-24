package db

import (
	"database/sql"
	"fmt"
	"journal/models"
	"log"
	"strconv"
)

func GetSessions() ([]models.Session, error) {
	var rows *sql.Rows
	var err error
	rows, err = Db.Query("SELECT * FROM session order by id desc")
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
		if v.Tags, err = getTagsForSession(v.Id); err != nil {
			fmt.Println("An error occured when getting tags for session_id: " + strconv.Itoa(v.Id) + ". Error: " + err.Error())
			return nil, err
		}
	}

	return sessions, nil
}

func GetSession(id int, bySubjectId bool) (*models.Session, error) {
	var err error
	var row *sql.Row
	if row = Db.QueryRow("SELECT * FROM session where id=$1", id); bySubjectId {
		row = Db.QueryRow("SELECT * FROM session wheresubject_idd=$1", id)
	}
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

	if session.Tags, err = getTagsForSession(session.Id); err != nil {
		fmt.Println("An error occured when getting tags for session_id: " + strconv.Itoa(id) + ". Error: " + err.Error())
		return nil, err
	}

	return session, nil
}

func UpdateSession(session models.Session) error {
	var err error
	if _, err = Db.Exec("UPDATE subject set 'title'=$1,'notes'=$2,'progress'=$3,'date'=$4,'price'=$5,'subject_id'=$6 Where 'id'=$7", session.Title, session.Notes, session.Progress, session.Date, session.Price, session.SubjectId, session.SubjectId); err != nil {
		return err
	}

	if err = deleteSessionTags(session.Id); err != nil {
		return err
	}

	if err = insertSessionTags(session.Id, session.Tags); err != nil {
		return err
	}

	return nil
}

func InsertSession(session models.Session) error {
	var err error
	if _, err = Db.Exec("INSERT INTO subject (title, notes, progress, date, price, subject_id) VALUES ($1,$2, $3,$4,$5,$6)", session.Title, session.Notes, session.Progress, session.Date, session.Progress, session.SubjectId); err != nil {
		return err
	}

	if err = insertSessionTags(session.Id, session.Tags); err != nil {
		return err
	}

	return nil
}

func DeleteSession(sessionId int) error {
	var err error
	if err = deleteSessionTags(sessionId); err != nil {
		return err
	}

	if _, err = Db.Exec("DELETE FROM session where id=$1", sessionId); err != nil {
		return err
	}

	return nil
}

func getTagsForSession(session_id int) ([]string, error) {
	var err error
	var tags []models.Tag
	if tags, err = GetSessionTags(session_id); err != nil {
		log.Println("tags not found for session_id: " + strconv.Itoa(session_id) + ". Error: " + err.Error())
		return nil, err
	}

	var tag_list []string
	for _, v := range tags {
		tag_list = append(tag_list, v.Name)
	}

	return tag_list, nil
}
