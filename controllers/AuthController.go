package controllers

import (
	"context"
	"core/database/postgres"
	"core/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)
type credentials struct {
	Login string `json:"login"`
	Password string `json:"password"`
}
func Authenticate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//clients []models.Client
	Client := models.Client{}

	var cred credentials

	db := postgres.Connect()

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&cred)

	err := db.QueryRow(context.Background(), postgres.LoginSQL, cred.Login).Scan(
		&Client.ID,
		&Client.Name,
		&Client.Surname,
		&Client.Login,
		&Client.Password,
		&Client.Age,
		&Client.Gender,
		&Client.Phone,
		&Client.Status,
		&Client.VerifiedAt)
	fmt.Println(Client)
	if err != nil {
		fmt.Println("error", err)
	}

	if MakeHash(cred.Password) == Client.Password {
		json.NewEncoder(w).Encode(Client)
	} else {
		w.WriteHeader(401)
	}
}

func MakeHash(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}