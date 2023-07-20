package main

import (
	"fmt"
	"traQuest/model"
	"traQuest/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	_, err = model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}

	router.SetRouting()
}
