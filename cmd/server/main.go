package main

import (
	"fmt"
	"github.com/erdemcemal/basket-service/internal/basket"
	"github.com/erdemcemal/basket-service/internal/database"
	basketstore "github.com/erdemcemal/basket-service/internal/store/basket"
	transportHttp "github.com/erdemcemal/basket-service/internal/transport/http"
	"log"
	"net/http"
)

type App struct {
}

func (a *App) Run() error {
	fmt.Println("setting up our application")

	var err error

	if err != nil {
		fmt.Println("failed to setup connection to the database")
		return err
	}
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bs := basketstore.NewStore(db)
	basketService := basket.NewService(bs)

	handler := transportHttp.NewHandler(basketService)

	if err := http.ListenAndServe(":3000", handler.Router); err != nil {
		fmt.Println("Failed to set up server")
		return err
	}
	fmt.Println("Server is running on port 3000")
	return nil
}

func main() {
	app := &App{}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
