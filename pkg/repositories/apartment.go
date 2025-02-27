package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/einherij/apt-manager/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Apartment struct {
	psql *sql.DB
}

func NewApartment(psql *sql.DB) *Apartment {
	return &Apartment{
		psql: psql,
	}
}

func (r *Apartment) All(ctx context.Context) (models.ApartmentSlice, error) {
	all, err := models.Apartments().All(ctx, r.psql)
	if err != nil {
		return nil, fmt.Errorf("error getting apartments: %w", err)
	}
	return all, nil
}

func (r *Apartment) Find(ctx context.Context, id int) (*models.Apartment, error) {
	apartment, err := models.FindApartment(ctx, r.psql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting apartment: %w", err)
	}
	return apartment, nil
}

func (r *Apartment) FindByBuildingID(ctx context.Context, buildingID int) (models.ApartmentSlice, error) {
	apartments, err := models.Apartments(
		models.ApartmentWhere.BuildingID.EQ(
			null.IntFrom(buildingID))).
		All(ctx, r.psql)
	if err != nil {
		return nil, fmt.Errorf("error getting apartments: %w", err)
	}
	return apartments, nil
}

func (r *Apartment) Upsert(ctx context.Context, apartment *models.Apartment) error {
	err := apartment.Upsert(
		ctx,
		r.psql,
		true,
		[]string{"id"},
		boil.Blacklist("id"),
		boil.Infer(),
	)
	if err != nil {
		return fmt.Errorf("error upserting apartment: %w", err)
	}
	return nil
}

func (r *Apartment) Delete(ctx context.Context, id int) error {
	apartment, err := models.FindApartment(ctx, r.psql, id)
	if err != nil {
		return fmt.Errorf("error getting apartment: %w", err)
	}

	_, err = apartment.Delete(ctx, r.psql)
	if err != nil {
		return fmt.Errorf("error deleting apartment: %w", err)
	}
	return nil
}
