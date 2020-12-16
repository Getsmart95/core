package postgres

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"log"
)

var db *pgxpool.Pool

func Connect()(db *pgxpool.Pool){
	db, err := pgxpool.Connect(context.Background(), `postgresql://root@localhost:5432/postgres?sslmode=disable`)

	if err != nil {
		log.Fatalf("Ошибка открытия базы данных %s", err)
	}
	return
}
