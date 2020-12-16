package models

import "time"

type Account struct {
	ID            int64
	ClientId      int64
	AccountNumber int64
	Balance       int64
	Status        bool
	CardNumber    string
	LimitTransfer int64
	LimitPayment  int64
	CreatedAt     time.Time
	UntilAt       time.Time
}

type AccountList struct {
	AccountWithUserName []AccountWithUserName
}

type AccountWithUserName struct {
	Account Account
	Client  Client
}

type AccountForUser struct {
	ID            int64
	ClientId      int64
	AccountNumber int64
	Balance       int64
	Status        bool
	CardNumber 	  string
}
