package models

type User struct {
	Id           string
	Name         string
	Email        string
	Password     string
	ScholarID    string
	Gender       string
	Hostel       string
	TokenVersion int
}
