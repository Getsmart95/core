package models

type Service struct{
	ID int
	Name string
	ServiceAccountNumber int64
}

type ServiceList struct{
	Services []Service
}
