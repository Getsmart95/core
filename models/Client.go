package models

import "time"

type Client struct {
	ID         int		`json:"ID"`
	Name       string 	`json:"Name"`
	Surname    string	`json:"Surname"`
	Login      string	`json:"Login"`
	Password   string 	`json:"Password"`
	Age        int		`json:"Age"`
	Gender     string	`json:"Gender"`
	Phone      string	`json:"Phone"`
	Status 	   bool		`json:"Status"`
	VerifiedAt time.Time`json:"VerifiedAt"`
}

type (
	ClientList struct {
		Clients []Client
	}
)
