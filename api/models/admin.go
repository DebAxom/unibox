package models

type Admin struct {
	Id           string
	Username     string
	Name         string
	Password     string
	Department   string
	TokenVersion int
}
