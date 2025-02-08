package main

import (
	"fmt"
	"log"
	"net/http"
	"warehouse-service/conn/mysql"
	"warehouse-service/conn/rabbitmq"
	productWarehouseHandler "warehouse-service/handler/product_warehouse"
	warehouseHandler "warehouse-service/handler/warehouse"
	"warehouse-service/middleware"
	productWarehouseRepo "warehouse-service/repository/product_warehouse"
	warehouseRepo "warehouse-service/repository/warehouse"
	productWarehouseUsecase "warehouse-service/usecase/product_warehouse"
	warehouseUsecase "warehouse-service/usecase/warehouse"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	rabbitmq.Connect()
	router := mux.NewRouter()

	rabbitPublisher := rabbitmq.NewRabbitPublisher(rabbitmq.RabbitConn)

	warehouseRepository := warehouseRepo.NewWarehouseRepository(mysql.MySQL)
	warehouseUsecase := warehouseUsecase.NewWarehouseUsecase(warehouseRepository)
	warehouseHandler := warehouseHandler.NewWarehouseHandler(warehouseUsecase)
	router.Handle("/warehouse/register", middleware.JWTMiddleware(http.HandlerFunc(warehouseHandler.Register))).Methods(http.MethodPost)
	router.Handle("/warehouse/update-status", middleware.JWTMiddleware(http.HandlerFunc(warehouseHandler.UpdateStatus))).Methods(http.MethodPut)

	productWarehouseRepository := productWarehouseRepo.NewProductWarehouseRepository(mysql.MySQL)
	productWarehouseUsecase := productWarehouseUsecase.NewProductWarehouseUsecase(productWarehouseRepository, rabbitPublisher, mysql.MySQL)
	productWarehouseHandler := productWarehouseHandler.NewProductWarehouseHandler(productWarehouseUsecase)
	router.Handle("/product-warehouse/register", middleware.JWTMiddleware(http.HandlerFunc(productWarehouseHandler.Register))).Methods(http.MethodPost)
	router.Handle("/product-warehouse/transfer", middleware.JWTMiddleware(http.HandlerFunc(productWarehouseHandler.TranserStockRequest))).Methods(http.MethodPost)
	router.Handle("/product-warehouse/add", middleware.JWTMiddleware(http.HandlerFunc(productWarehouseHandler.AddStockRequest))).Methods(http.MethodPost)
	router.Handle("/product-warehouse/deduct", middleware.JWTMiddleware(http.HandlerFunc(productWarehouseHandler.DeductStockRequest))).Methods(http.MethodPost)

	rabbitConsumer := rabbitmq.NewRabbitConsumer(rabbitmq.RabbitConn, productWarehouseHandler)
	go rabbitConsumer.ConsumeEvents()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8003", router)
	if err != nil {
		log.Fatal(err)
	}
}
