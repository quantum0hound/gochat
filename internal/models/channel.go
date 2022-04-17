package models

type Channel struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Creator     int    `json:"creator" db:"creator"`
}
