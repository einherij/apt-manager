package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/einherij/apt-manager/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Building struct {
	psql *sql.DB
}

func NewBuilding(psql *sql.DB) *Building {
	return &Building{
		psql: psql,
	}
}

func (r *Building) All(ctx context.Context) (models.BuildingSlice, error) {
	all, err := models.Buildings().All(ctx, r.psql)
	if err != nil {
		return nil, fmt.Errorf("error getting buildings: %w", err)
	}
	return all, nil
}

func (r *Building) Find(ctx context.Context, id int) (*models.Building, error) {
	building, err := models.FindBuilding(ctx, r.psql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting building: %w", err)
	}
	return building, nil
}

func (r *Building) Upsert(ctx context.Context, building *models.Building) error {
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	err = building.Upsert(
		ctx,
		tx,
		true,
		[]string{"id"},
		boil.Blacklist("id"),
		boil.Infer(),
	)
	if err != nil {
		return fmt.Errorf("error upserting building: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (r *Building) Delete(ctx context.Context, id int) error {
	tx, err := r.psql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	building, err := models.FindBuilding(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("error getting building: %w", err)
	}
	_, err = building.Delete(ctx, tx)
	if err != nil {
		return fmt.Errorf("error deleting building: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}
