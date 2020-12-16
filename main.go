package main

import (
	"core/controllers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main(){

	router := httprouter.New()

	router.POST("/api/login", controllers.Authenticate)
	router.GET("/api/getAllAccounts/:id", controllers.GetAccountsByIdHandler)
	router.POST("/api/transferTo/:cardNumber", controllers.TransferToByHandler) //Оптимизировать код
	router.GET("/api/getAllServices", controllers.GetAllServices)
	router.POST("/api/payService", controllers.PayServiceByHandler)
	err := http.ListenAndServe(":8088", router)

	if err != nil{
		log.Fatal("ListenAndServe:", err)
	}
}
