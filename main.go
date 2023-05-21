package main

import (
	"fmt"
	"traQuest/model"
	"traQuest/router"

	"github.com/joho/godotenv"
	"github.com/srinathgs/mysqlstore"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	db, err := model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 3600, []byte("secret-token"))
	if err != nil {
		panic(err)
	}
	router.SetRouting(store)
}
