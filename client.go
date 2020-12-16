package main

import (
	"onlineBanking/cli/constants"
	"onlineBanking/cli/controllers"
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"log"
	"os"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), `postgresql://root@localhost:5432/postgres?sslmode=disable`)

	if err != nil {
		log.Fatalf("Ошибка открытия базы данных %s", err)
	}
	defer db.Close()

	AutorizeByClient(db)
}
func AutorizeByClient(db *pgxpool.Pool) {
	var cmd string
	for {
		fmt.Println(constants.UnauthorizedOperations)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			id, err := controllers.Authorize(db)
			if err != nil {
				fmt.Println("Попробуйте еще раз")
				continue
			} else {
				controllers.AuthorizedOperations(id, db)
			}
		case "2":
			err := controllers.GetATMsForClient(db)
			if err != nil {
				fmt.Printf("Ошибка выдачи списка банкоматов %s:", err)
			}
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Введена неверная команда, попробуйте еще раз", constants.UnauthorizedOperations)
			continue
		}
	}
}
