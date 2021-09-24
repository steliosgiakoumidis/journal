package db

import (
	"database/sql"
	"errors"
	"fmt"
	"journal/models"
)

func GetAllTags() ([]models.Tag, error){
	var rows *sql.Rows
	var err error
	rows, err = Db.Query("SELECT * FROM tag order by name")
	if err != nil{
		fmt.Println("An error occured when getting all subjects. Error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil{
			println("An error occured when assigning tags. Error: " + err.Error())
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func InsertTag(tag models.Tag) error {
	var err error
	if _, err = Db.Exec("INSERT INTO tag (Name) VALUES ($1)", tag.Name); err != nil {
		return err
	}

	return nil
}

func DeleteTag(id int) error{
	var err error
	if _, err = Db.Exec("DELETE from session_tag where tag_id=$1", id); err != nil{
		return err
	}
	if _, err = Db.Exec("DELETE from tag where id=$1", id); err != nil{
		return err
	}

	return nil
}

func GetSessionTags(sessionId int) ([]models.Tag, error) {
	var err error
	var rows *sql.Rows
	rows, err = Db.Query(
		"SELECT t.id, t.name FROM tag t inner join session_tag st on st.tag_id = t.id where st.session_id = $1", sessionId)
	if err != nil{
		fmt.Println("An error occured when getting all subjects. Error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil{
			println("An error occured when assigning tags. Error: " + err.Error())
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func deleteSessionTags(session_id int) error{
	var err error
	if _, err = Db.Exec("DELETE from session_tag where session_id=$1", session_id); err != nil{
		return err
	}

	return nil
}

func insertSessionTags(session_id int, tags_list []string) error{
	var tags []models.Tag
	var err error
	if tags, err = GetAllTags(); err != nil{
		return err
	}

	tag_ids_list := []int{}
	for i, v := range tags_list{
		if tag_ids_list[i], err = getTagId(tags, v); err != nil{
			continue
		}
	}

	if err = insertTags(session_id, tag_ids_list); err != nil{
		return err
	}

	return nil
}

func insertTags(session_id int, tag_ids_list []int) error {
	var err error
	for _, v := range tag_ids_list {
		if _, err = Db.Exec("INSERT INTO session_tag (session_id, tag_id) values ($1, $2)", session_id, v); err != nil{
			return err
		}
	}

	return nil
}

func getTagId(tags []models.Tag, tag_name string) (int, error) {
	for _, v := range tags{
		if v.Name == tag_name{
			return v.Id, nil
		}
	}

	return 0, errors.New("Tag not found for tag name: " + tag_name)
}