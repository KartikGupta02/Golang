package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/controller"
	"main/route"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:K@rtik3275@(127.0.0.1:3306)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Connected with database")

	defer db.Close()

	controller.SetDB(db)

	route.RegisterUserRoutes()
	route.RegisterVideoRoutes()

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
