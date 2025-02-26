package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/einherij/apt-manager/models"
	"github.com/einherij/apt-manager/pkg/repositories"
)

type BuildingHandler struct {
	repository BuildingRepository
}

type BuildingRepository interface {
	All(ctx context.Context) (models.BuildingSlice, error)
	Find(ctx context.Context, id int) (*models.Building, error)
	Upsert(ctx context.Context, building *models.Building) error
	Delete(ctx context.Context, id int) error
}

var _ BuildingRepository = new(repositories.Building)

func NewBuildingHandler(buildingRepository BuildingRepository) *BuildingHandler {
	return &BuildingHandler{
		repository: buildingRepository,
	}
}

func (r *BuildingHandler) All(c *fiber.Ctx) error {
	all, err := r.repository.All(c.Context())
	if err != nil {
		return fmt.Errorf("error getting all buildings: %w", err)
	}
	err = c.JSON(all)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}

func (r *BuildingHandler) Find(c *fiber.Ctx) error {
	buildingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("error parsing building ID: %w", err)
	}
	building, err := r.repository.Find(c.Context(), buildingID)
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
}

func (r *BuildingHandler) Upsert(c *fiber.Ctx) error {
	var building = new(models.Building)
	if err := c.BodyParser(building); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}
	if err := r.repository.Upsert(c.Context(), building); err != nil {
		return fmt.Errorf("error upserting building: %w", err)
	}
	return nil
}

func (r *BuildingHandler) Delete(c *fiber.Ctx) error {
	buildingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("error parsing building ID: %w", err)
	}
	if err := r.repository.Delete(c.Context(), buildingID); err != nil {
		return fmt.Errorf("error deleting building: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return nil
}
