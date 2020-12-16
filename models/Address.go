package models

type Address struct {
	ID       int
	Country  string
	City     string
	Street   string
	home     int
}

type AddressList struct {
	Addresses []Address
}
