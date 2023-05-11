package main

import (
	"fmt"
	"traQuest/model"
	"traQuest/router"
)

func main() {
	_, err := model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}
	router.SetRouting()
}
