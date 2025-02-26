package routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, buildingHandler *BuildingHandler, apartmentHandler *ApartmentHandler) {
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
