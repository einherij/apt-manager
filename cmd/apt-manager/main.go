package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/einherij/apt-manager/models"
	"github.com/einherij/apt-manager/pkg/repositories"
)

const (
	PG_USER     = "apt-manager-user"
	PG_PASSWORD = "pg-pass"
	PG_HOST     = "localhost"
	PG_PORT     = "5433"
	PG_DB       = "apt-manager"

	SERVER_ADDRESS = "127.0.0.1:8080"
)

func main() {
	var psqlConnection = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		PG_USER, PG_PASSWORD, PG_HOST, PG_PORT, PG_DB)

	m := Must(migrate.New(
		"file://migrations",
		psqlConnection,
	))
	err := m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		PanicOnError(err)
	}

	psql := Must(sql.Open("postgres", psqlConnection))
	defer func() {
		if err = psql.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()
	buildingRepository := repositories.New(psql)
	apartmentRepository := repositories.NewApartment(psql)

	app := fiber.New()

	registerBuildingRoutes(app, buildingRepository)
	registerApartmentRoutes(app, apartmentRepository)

	log.Fatal(app.Listen(SERVER_ADDRESS))
}

func registerBuildingRoutes(app *fiber.App, buildingRepository *repositories.Building) {
	app.Get("/buildings", func(c *fiber.Ctx) error {
		all, err := buildingRepository.All(c.Context())
		if err != nil {
			return fmt.Errorf("error getting all buildings: %w", err)
		}
		err = c.JSON(all)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
		return nil
	})

	app.Get("/buildings/:id", func(c *fiber.Ctx) error {
		buildingID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fmt.Errorf("error parsing building ID: %w", err)
		}
		building, err := buildingRepository.Find(c.Context(), buildingID)
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err != nil {
			return fmt.Errorf("error finding building: %w", err)
		}
		err = c.JSON(building)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
		return nil
	})
	app.Post("/buildings", func(c *fiber.Ctx) error {
		var building = new(models.Building)
		if err := c.BodyParser(building); err != nil {
			return fmt.Errorf("error parsing request body: %w", err)
		}
		if err := buildingRepository.Upsert(c.Context(), building); err != nil {
			return fmt.Errorf("error upserting building: %w", err)
		}
		return nil
	})
	app.Delete("/buildings/:id", func(c *fiber.Ctx) error {
		buildingID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fmt.Errorf("error parsing building ID: %w", err)
		}
		if err := buildingRepository.Delete(c.Context(), buildingID); err != nil {
			return fmt.Errorf("error deleting building: %w", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return nil
	})
}

func registerApartmentRoutes(app *fiber.App, apartmentRepository *repositories.Apartment) {
	app.Get("/apartments", func(c *fiber.Ctx) error {
		all, err := apartmentRepository.All(c.Context())
		if err != nil {
			return fmt.Errorf("error getting all apartments: %w", err)
		}
		err = c.JSON(all)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
		return nil
	})
	app.Get("/apartments/:id", func(c *fiber.Ctx) error {
		apartmentID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fmt.Errorf("error parsing apartment ID: %w", err)
		}
		apartment, err := apartmentRepository.Find(c.Context(), apartmentID)
		err = c.JSON(apartment)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
		return nil
	})
	app.Get("/apartments/building/:buildingId", func(c *fiber.Ctx) error {
		buildingID, err := strconv.Atoi(c.Params("buildingId"))
		if err != nil {
			return fmt.Errorf("error parsing building ID: %w", err)
		}
		apartments, err := apartmentRepository.FindByBuildingID(c.Context(), buildingID)
		err = c.JSON(apartments)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
		return nil
	})
	app.Post("/apartments", func(c *fiber.Ctx) error {
		var apartment = new(models.Apartment)
		if err := c.BodyParser(apartment); err != nil {
			return fmt.Errorf("error parsing request body: %w", err)
		}
		if err := apartmentRepository.Upsert(c.Context(), apartment); err != nil {
			return fmt.Errorf("error upserting apartment: %w", err)
		}
		return nil
	})
	app.Delete("/apartments/:id", func(c *fiber.Ctx) error {
		apartmentID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fmt.Errorf("error parsing apartment ID: %w", err)
		}
		if err := apartmentRepository.Delete(c.Context(), apartmentID); err != nil {
			return fmt.Errorf("error deleting apartment: %w", err)
		}
		return nil
	})
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
