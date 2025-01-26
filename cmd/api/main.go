package main

import (
	"github.com/ClearingHouse/internal/app"
	"log"
)

func main() {
	db, err := app.InitDataBase()
	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDB()

	app := app.NewApp(db)
	if err := app.Run(); err != nil {
		log.Fatal("Server encountered an error:", err)
	}

}
