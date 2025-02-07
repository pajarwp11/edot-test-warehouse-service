package main

import (
	"fmt"
	"log"
	"net/http"
	"warehouse-service/db/mysql"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	router := mux.NewRouter()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8003", router)
	if err != nil {
		log.Fatal(err)
	}
}
