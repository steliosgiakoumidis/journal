package db

import (
	"database/sql"
	"fmt"
	"journal/models"
	"strconv"
)

type SubjectRepository struct{}

func (s SubjectRepository) GetSubjects() ([]models.Subject, error) {
	var rows *sql.Rows
	var err error
	rows, err = DbConn.Query("SELECT * FROM subject order by firstname")
	if err != nil {
		fmt.Println("An error occured when getting all subjects. Error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var subjects = []models.Subject{}
	for rows.Next() {
		subject := models.Subject{}
		if err = rows.Scan(&subject.Id, &subject.Firstname, &subject.LastName, &subject.Phonenumber, &subject.Email, &subject.Age, &subject.AgreedPrice); err != nil {
			println("An error occured when assigning subjects. Error: " + err.Error())
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (s SubjectRepository) GetSubject(id int) (*models.Subject, error) {
	var err error
	row := DbConn.QueryRow("SELECT * FROM SUBJECT WHERE id=$1", id)
	if row.Err() != nil {
		println("An error occured when getting subject with lastname " + strconv.Itoa(id) + ". Error: " + err.Error())
		return nil, err
	}
	var subject *models.Subject
	err = row.Scan(&subject.Id, &subject.Firstname, &subject.LastName, &subject.Phonenumber, &subject.Email, &subject.Age, &subject.AgreedPrice)
	switch err {
	case sql.ErrNoRows:
		println("Subject not found with lastname " + strconv.Itoa(id) + ". Error: " + err.Error())
		return nil, err
	case nil:
		println("subject found")
		return subject, nil
	default:
		panic("What the fuck")
	}
}

func (s SubjectRepository) UpdateSubject(subject models.Subject) error {
	var err error
	if _, err = DbConn.Exec("UPDATE subject set 'firstname'=$1,'lastname'=$2,'phonenumber'=$3,'email'=$4,'age'=$5,'agreedprice'=$6 Where 'id'=$7", subject.Firstname, subject.LastName, subject.Phonenumber, subject.Email, subject.Age, subject.AgreedPrice, subject.Id); err != nil {
		return err
	}

	return nil
}

func (s SubjectRepository) InsertSubject(subject models.Subject) error {
	var err error
	if _, err = DbConn.Exec("INSERT INTO subject (Firstname, LastName, Phonenumber, Email, Age, AgreedPrice) VALUES ($1,$2, $3,$4,$5,$6)", subject.Firstname, subject.LastName, subject.Phonenumber, subject.Email, subject.Age, subject.AgreedPrice); err != nil {
		return err
	}

	return nil
}

func (s SubjectRepository) DeleteSubject(id int) error {
	var err error
	if _, err = DbConn.Exec("DELETE from subject where id=$1", id); err != nil {
		return err
	}

	return nil
}
