package main

import (
	"core/controllers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main(){

	router := httprouter.New()

	router.POST("/login", controllers.Authenticate)
	router.GET("/getAllAccounts/:id", controllers.GetAccountsByIdHandler)
	router.POST("/transferTo/:cardNumber", controllers.TransferToByHandler)
	err := http.ListenAndServe(":8088", router)

	if err != nil{
		log.Fatal("ListenAndServe:", err)
	}
}
