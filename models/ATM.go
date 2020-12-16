package models

type ATM struct {
	ID        int
	Name      string
	Status    bool
	//CreatedAt time.Time
}

type ATMList struct {
	ATMs []ATM
}
