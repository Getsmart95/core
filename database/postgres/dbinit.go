package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
)

func DbInit(database *pgxpool.Pool) error {
	var DDLs = []string{clientsTable, clientsAccountsTable, ATMsTable, servicesTable, historiesTable, atmAddresses}
	for _, ddl := range DDLs {
		_, err := database.Exec(context.Background(), ddl)
		if err != nil {
			fmt.Println("Can't init this #{ddl} error is ${err}")
		}
	}
	return nil
}
