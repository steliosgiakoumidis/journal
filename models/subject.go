package models

type Subject struct {
	Id int `json:"id"`
	Firstname string `json:"firstname"`
	LastName string `json:"lastname"`
	Phonenumber string `json:"phonenumber"`
	Email string `json:"email"`
	Age int `json:"age"`
	AgreedPrice int `json:"agreedprice"`

	//TODO add createdAt updatedAt
}

func (s *Subject) IsValid() bool {
	if s.Firstname == ""{
		return false
	}
	if s.LastName == ""{
		return false
	}
	if s.Phonenumber == ""{
		return false
	}
	if s.Email == ""{
		return false
	}
	if s.Age == 0{
		return false
	}
	if s.AgreedPrice == 0{
		return false
	}

	return true
}
