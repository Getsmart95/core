package models

type Service struct{
	ID int
	Name string
	AccountNumber int64
}

type ServiceList struct{
	Services []Service
}
