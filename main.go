package main

import (
	"fmt"
	"log"
	"net/http"
	"warehouse-service/conn/mysql"
	"warehouse-service/conn/rabbitmq"
	warehouseHandler "warehouse-service/handler/warehouse"
	"warehouse-service/middleware"
	warehouseRepo "warehouse-service/repository/warehouse"
	warehouseUsecase "warehouse-service/usecase/warehouse"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	rabbitmq.Connect()
	router := mux.NewRouter()
	warehouseRepository := warehouseRepo.NewWarehouseRepository(mysql.MySQL)
	warehouseUsecase := warehouseUsecase.NewWarehouseUsecase(warehouseRepository)
	warehouseHandler := warehouseHandler.NewWarehouseHandler(warehouseUsecase)
	router.Handle("/warehouse/register", middleware.JWTMiddleware(http.HandlerFunc(warehouseHandler.Register))).Methods(http.MethodPost)

	fmt.Println("server is running")
	err := http.ListenAndServe(":8003", router)
	if err != nil {
		log.Fatal(err)
	}
}
