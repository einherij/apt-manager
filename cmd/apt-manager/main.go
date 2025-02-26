package main

import (
	"database/sql"
	"log"

	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/einherij/apt-manager/pkg/config"
	"github.com/einherij/apt-manager/pkg/repositories"
	"github.com/einherij/apt-manager/pkg/routes"
)

func main() {
	var conf = config.NewConfig()
	PanicOnError(conf.ParseEnv())

	m := Must(migrate.New(
		"file://migrations",
		conf.Postgres.PostgresConnection(),
	))
	err := m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		PanicOnError(err)
	}

	psql := Must(sql.Open("postgres", conf.Postgres.PostgresConnection()))
	defer func() {
		if err = psql.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()
	buildingRepository := repositories.NewBuilding(psql)
	apartmentRepository := repositories.NewApartment(psql)

	app := fiber.New()

	buildingHandler := routes.NewBuildingHandler(buildingRepository)
	apartmentHandler := routes.NewApartmentHandler(apartmentRepository)
	registerRoutes(app, buildingHandler, apartmentHandler)

	err = app.Listen(conf.ServerAddress)
	log.Fatal(app.Listen(conf.ServerAddress))
}

func registerRoutes(app *fiber.App, buildingHandler *routes.BuildingHandler, apartmentHandler *routes.ApartmentHandler) {
	app.Get("/apartments", apartmentHandler.All)
	app.Get("/apartments/:id", apartmentHandler.Find)
	app.Get("/apartments/building/:buildingId", apartmentHandler.FindByBuildingID)
	app.Post("/apartments", apartmentHandler.Upsert)
	app.Delete("/apartments/:id", apartmentHandler.Delete)

	app.Get("/buildings", buildingHandler.All)
	app.Get("/buildings/:id", buildingHandler.Find)
	app.Post("/buildings", buildingHandler.Upsert)
	app.Delete("/buildings/:id", buildingHandler.Delete)
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
