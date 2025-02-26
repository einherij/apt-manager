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
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	all, err := models.Apartments().All(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("error getting apartments: %w", err)
	}
	return all, nil
}

func (r *Apartment) Find(ctx context.Context, id int) (*models.Apartment, error) {
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	apartment, err := models.FindApartment(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting apartment: %w", err)
	}
	return apartment, nil
}

func (r *Apartment) FindByBuildingID(ctx context.Context, buildingID int) (models.ApartmentSlice, error) {
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	apartments, err := models.Apartments(
		models.ApartmentWhere.BuildingID.EQ(
			null.NewInt(buildingID, true))).
		All(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("error getting apartments: %w", err)
	}
	return apartments, nil
}

func (r *Apartment) Upsert(ctx context.Context, apartment *models.Apartment) error {
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	err = apartment.Upsert(
		ctx,
		tx,
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
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	apartment, err := models.FindApartment(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("error getting apartment: %w", err)
	}
	_, err = apartment.Delete(ctx, tx)
	if err != nil {
		return fmt.Errorf("error deleting apartment: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}
