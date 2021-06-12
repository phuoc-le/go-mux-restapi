package main

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_NAME"))

	// port = strconv.Atoi(os.Getenv("PORT"))

	fmt.Println("Start Api with Port: ", os.Getenv("PORT"))

	a.Run(":"+ os.Getenv("PORT"))
}
