package routes

import (
	"context"
	"fmt"
	"github.com/einherij/apt-manager/models"
	"github.com/einherij/apt-manager/pkg/repositories"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ApartmentRepository interface {
	All(ctx context.Context) (models.ApartmentSlice, error)
	Find(ctx context.Context, id int) (*models.Apartment, error)
	FindByBuildingID(ctx context.Context, buildingID int) (models.ApartmentSlice, error)
	Upsert(ctx context.Context, apartment *models.Apartment) error
	Delete(ctx context.Context, id int) error
}

var _ ApartmentRepository = new(repositories.Apartment)

type ApartmentHandler struct {
	repository ApartmentRepository
}

func NewApartmentHandler(apartmentRepository ApartmentRepository) *ApartmentHandler {
	return &ApartmentHandler{
		repository: apartmentRepository,
	}
}

func (r *ApartmentHandler) All(c *fiber.Ctx) error {
	all, err := r.repository.All(c.Context())
	if err != nil {
		return fmt.Errorf("error getting all apartments: %w", err)
	}
	err = c.JSON(all)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}

func (r *ApartmentHandler) Find(c *fiber.Ctx) error {
	apartmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("error parsing apartment ID: %w", err)
	}
	apartment, err := r.repository.Find(c.Context(), apartmentID)
	err = c.JSON(apartment)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}

func (r *ApartmentHandler) FindByBuildingID(c *fiber.Ctx) error {
	buildingID, err := strconv.Atoi(c.Params("buildingId"))
	if err != nil {
		return fmt.Errorf("error parsing building ID: %w", err)
	}
	apartments, err := r.repository.FindByBuildingID(c.Context(), buildingID)
	err = c.JSON(apartments)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}

func (r *ApartmentHandler) Upsert(c *fiber.Ctx) error {
	var apartment = new(models.Apartment)
	if err := c.BodyParser(apartment); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}
	if err := r.repository.Upsert(c.Context(), apartment); err != nil {
		return fmt.Errorf("error upserting apartment: %w", err)
	}
	return nil
}

func (r *ApartmentHandler) Delete(c *fiber.Ctx) error {
	apartmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("error parsing apartment ID: %w", err)
	}
	if err := r.repository.Delete(c.Context(), apartmentID); err != nil {
		return fmt.Errorf("error deleting apartment: %w", err)
	}
	return nil
}
