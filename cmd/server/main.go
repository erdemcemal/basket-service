package main

import (
	"github.com/erdemcemal/basket-service/internal/basket"
	"github.com/erdemcemal/basket-service/internal/database"
	basketstore "github.com/erdemcemal/basket-service/internal/store/basket"
	transportHttp "github.com/erdemcemal/basket-service/internal/transport/http"
	log "github.com/siruspen/logrus"
)

// App - contains the application configuration.
type App struct {
	Name    string
	Version string
}

// Run - sets up our application and starts the server.
func (a *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    a.Name,
			"AppVersion": a.Version,
		},
	).Info("Setting up application")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		log.Error(err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		log.Error(err)
		return err
	}
	bs := basketstore.NewBasketStore(db)
	basketService := basket.NewService(bs)

	handler := transportHttp.NewHandler(basketService)
	if err := handler.Serve(); err != nil {
		log.Error("Failed to set up server")
		return err
	}
	return nil
}

func main() {
	app := &App{
		Name:    "basket-service",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error("Failed to run application")
		log.Fatal(err)
	}
}
