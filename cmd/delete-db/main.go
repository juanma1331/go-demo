package main

import (
	"fmt"
	"go-demo/internal/infra"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db, err := infra.OpenDB()

	db.Exec("DROP TABLE IF EXISTS users;")
	db.Exec("DROP TABLE IF EXISTS products;")
	db.Exec("DROP TABLE IF EXISTS carts;")
	db.Exec("DROP TABLE IF EXISTS cart_details;")
	db.Exec("DROP TABLE IF EXISTS http_sessions;")

	if err != nil {
		panic(err)
	}

}
