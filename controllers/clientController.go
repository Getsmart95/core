package controllers

import (
	"context"
	"core/database/postgres"
	"core/models"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"onlineBanking/core/packages"
	"strconv"
)

type clientAccount struct {
	AccountNumber int64 `json:"accountNumber"`
	TransferCardNumber string `json:"transferCardNumber"`
	TransferAmmount string `json:"transferAmmount"`
	Message string `json:"message"`
}

type service struct {
	Id int64 `json:"ID"`
	Name string `json:"name"`
	AccountNumber int64 `json:"accountNumber"`
	Ammount string `json:"ammount"`
	ServiceAccountNumber int64 `json:"serviceAccountNumber"`
}

func GetATMsForClient(db *pgxpool.Pool) (err error) {
	ms, err := services.GetAllATMs(db)
	if err != nil {
		return err
	}
	i := 0
	for _, value := range ms {
		i++
		fmt.Println(value)
	}
	if i == 0 {
		fmt.Println("Список банкоматов пуст")
	}
	return nil
}

//go install ./...
func GetAccountsByIdHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db := postgres.Connect()
	//fmt.Fprint(w, "hello")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(p.ByName("id"))
	//list, err := SearchAccountById(id, db)

	rows, err := db.Query(context.Background(), postgres.SearchAccountByID, id)

	if err != nil {
		fmt.Errorf("Активных аккаунтов нет %e\n", err)
	}
	Accounts := []models.AccountForUser{}
	Account := models.AccountForUser{}

	for rows.Next() {
		rows.Scan(
			&Account.ID,
			&Account.ClientId,
			&Account.AccountNumber,
			&Account.Balance,
			&Account.Status,
			&Account.CardNumber)

		Accounts = append(Accounts, Account)
	}

	json.NewEncoder(w).Encode(Accounts)

}

func TransferToByHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db := postgres.Connect()
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var data clientAccount
	TransferCardNumber := p.ByName("cardNumber")
	json.NewDecoder(r.Body).Decode(&data)
	AccountNumber := data.AccountNumber
	Amount, _ := strconv.Atoi(data.TransferAmmount)
	Message := data.Message
	fmt.Println(TransferCardNumber)
	fmt.Println(data)
	var TransferAccountNumber int64
	err := db.QueryRow(context.Background(), `select account_number from accounts where card_number = ($1)`, TransferCardNumber).Scan(&TransferAccountNumber)
	if err != nil {
		w.WriteHeader(400)
	}
	err = TransferToAccount(AccountNumber, TransferCardNumber, Amount, Message, TransferAccountNumber, db)
	if err != nil {
		w.WriteHeader(400)
	}
}

func TransferToAccount(AccountNumber int64, TransferCardNumber string, Amount int, Message string, TransferAccountNumber int64, db *pgxpool.Pool) (err error) {

	var ServiceId int64
	ServiceId = 1

	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance - ($1)
								                 where account_Number = ($2)`, Amount, AccountNumber)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance + ($1)
                                                where card_number = ($2)`, Amount, TransferCardNumber)

	if err != nil {
		return err
	}
	fmt.Println(AccountNumber,TransferCardNumber,Amount, Message,ServiceId)
	_, err = db.Exec(context.Background(), `insert into histories(sender_account_number, recipient_account_number, money, message, service_id)
											values( $1, $2, $3, $4, $5 )`, AccountNumber, TransferAccountNumber, Amount, Message, ServiceId)
	if err != nil {
		return err
	}
	fmt.Println("Перевод денег успешно выполнено!")
	return nil
}

func GetAllServices(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db := postgres.Connect()
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	Services := []models.Service{}

	rows, err := db.Query(context.Background(), postgres.GetAllServices)
	if err != nil {
		w.WriteHeader(400)
	}

	for rows.Next(){
		Service := models.Service{}
		err := rows.Scan(&Service.ID, &Service.Name, &Service.ServiceAccountNumber)
		if err != nil {
			w.WriteHeader(400)
		}
		Services = append(Services, Service)
	}
	if rows.Err() != nil{
		w.WriteHeader(400)
	}
	json.NewEncoder(w).Encode(Services)

}

func PayServiceByHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db := postgres.Connect()
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var data service
	Services := []models.Service{}

	json.NewDecoder(r.Body).Decode(&data)
	Ammount, _ := strconv.Atoi(data.Ammount)
	err := Transfer(data.AccountNumber, Ammount, data.Id, db)
	if err != nil {
		fmt.Println("Перевод невозможен")
	}
	json.NewEncoder(w).Encode(Services)

}


func Transfer(accountNumber int64, Ammount int, ServiceID int64, db *pgxpool.Pool) (err error) {

	var Message string
	Message = "Оплата услуги"

	var ServiceAccountNumber int64
	err = db.QueryRow(context.Background(), `select account_number from services where id = ($1)`, ServiceID).Scan(&ServiceAccountNumber)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance - ($1) where account_number = ($2)`, Ammount, accountNumber)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance + ($1) where account_number = ($2)`, Ammount, ServiceAccountNumber)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), `insert into histories(sender_account_number, recipient_account_number, money, message, service_id)
											values( $1, $2, $3, $4, $5 )`, accountNumber, ServiceAccountNumber, Ammount, Message, ServiceID)
	if err != nil {
		return err
	}


	return nil
}





























//////
///////////////////
//func ChooseAccount(id int64, db *pgxpool.Pool) (AccountNumber int64, err error) {
//	fmt.Println("Выберите счет:")
//	accounts, err := SearchAccountByIdHandler(id, db)
//	if err != nil {
//		return -1, err
//	}
//	//	fmt.Println(accounts)
//
//	for {
//		var cmd int64
//		fmt.Scan(&cmd)
//		switch int64(len(accounts)) >= cmd && cmd > 0 {
//		case true:
//			return accounts[cmd], nil
//		case false:
//			fmt.Println("Введите заново в пределах количество ваших счетов")
//		}
//	}
//	return -1, nil
//}
//
/////////////////////////
////
//func PayServiceHandler(id int64, db *pgxpool.Pool) (err error) {
//	fmt.Println("Выберите счет:")
//	accounts, err := SearchAccountByIdHandler(id, db)
//	if err != nil {
//		return err
//	}
//
//	for {
//		var cmd int64
//		fmt.Scan(&cmd)
//		switch int64(len(accounts)) >= cmd && cmd > 0 {
//		case true:
//			ChooseToService(accounts[cmd], db)
//			return nil
//		case false:
//			fmt.Println("Введите заново в пределах количество ваших счетов")
//		}
//	}
//	return nil
//}
//
//func GetAllServicesHandler(db *pgxpool.Pool) (err error) {
//services, err := services.GetAllServices(db)
//if err != nil {
//	fmt.Errorf("Get all services didn't work %e", err)
//	return nil
//}
//
//for _, service := range services {
//	fmt.Println(service.ID, service.Name, service.AccountNumber)
//}
//return nil
//}
//
//func ChooseToService(AccountNumber int64, db *pgxpool.Pool) (err error) {
//	fmt.Println("Выберите услугу: ")
//	err = GetAllServicesHandler(db)
//	if err != nil {
//		fmt.Errorf("GetServiceHandler %e", err)
//		return err
//	}
//	for {
//		var cmd int64
//		fmt.Scan(&cmd)
//		err := services.CheckServiceHaving(cmd, db)
//		if err != nil {
//			fmt.Println("Такой услуги нет, попробуйте еще раз")
//			continue
//		} else {
//			fmt.Println("Введите сумму оплаты: ")
//			var Ammount int64
//			fmt.Scan(&Ammount)
//			err := Transfer(AccountNumber, Ammount, cmd, db)
//			if err != nil {
//				fmt.Println("Перевод невозможен")
//			}
//		}
//		return nil
//	}
//}
//
//func Transfer(accountNumber int64, Ammount int64, ServiceID int64, db *pgxpool.Pool) (err error) {
//
//	var Message string
//	Message = "Оплата услуги"
//
//	var ServiceAccountNumber int64
//	err = db.QueryRow(context.Background(), `select account_number from services where id = ($1)`, ServiceID).Scan(&ServiceAccountNumber)
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance - ($1) where account_number = ($2)`, Ammount, accountNumber)
//	if err != nil {
//		return err
//	}
//	_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance + ($1) where account_number = ($2)`, Ammount, ServiceAccountNumber)
//	if err != nil {
//		return err
//	}
//	_, err = db.Exec(context.Background(), `insert into histories(sender_account_number, recipient_account_number, money, message, service_id)
//											values( $1, $2, $3, $4, $5 )`, accountNumber, ServiceAccountNumber, Ammount, Message, ServiceID)
//	if err != nil {
//		return err
//	}
//	//_, err = db.Exec(context.Background(), `UPDATE accounts set balance = balance - ? where accountNumber = ?`, ServicePrice, accountNumber)
//	//if err != nil {
//	//	return err
//	//}
//
//	return nil
//}
