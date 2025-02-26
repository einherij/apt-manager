package repositories

import (
	"context"
	"github.com/einherij/apt-manager/models"
	"github.com/volatiletech/null/v8"
	"time"
)

func (s *RepositoriesSuite) TestApartmentRepository() {
	repo := NewApartment(s.psql)
	const testID = 1

	building := &models.Building{
		ID:      testID,
		Name:    null.StringFrom("TestBuildingName"),
		Address: null.StringFrom("TestBuildingAddress"),
	}
	apartment := &models.Apartment{
		ID:         testID,
		Number:     null.StringFrom("10"),
		BuildingID: null.IntFrom(testID),
		Floor:      null.IntFrom(4),
		SQMeters:   null.IntFrom(30),
	}
	_, _ = s.psql.Exec("INSERT INTO building (id, name, address) VALUES ($1, $2, $3)",
		building.ID, building.Name, building.Address)
	_, _ = s.psql.Exec("INSERT INTO apartment (id, number, building_id, floor, sq_meters) VALUES ($1, $2, $3, $4, $5)",
		apartment.ID, apartment.Number, apartment.BuildingID, apartment.Floor, apartment.SQMeters)

	s.Run("Upsert", func() {
		apartment.Floor = null.IntFrom(5)
		err := repo.Upsert(context.Background(), apartment)
		s.NoError(err)
		apt, err := repo.Find(context.Background(), testID)
		s.NoError(err)
		s.Equal(apartment.Floor, apt.Floor)
	})

	s.Run("All", func() {
		apartments, err := repo.All(context.Background())
		s.NoError(err)
		s.NotEmpty(apartments)
	})

	s.Run("Find", func() {
		apt, err := repo.Find(context.Background(), testID)
		s.NoError(err)
		s.NotNil(apt)
	})

	s.Run("FindByBuildingID", func() {
		apartments, err := repo.FindByBuildingID(context.Background(), testID)
		s.NoError(err)
		s.NotEmpty(apartments)
	})

	s.NoError(repo.Delete(context.Background(), testID))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.psql.ExecContext(ctx, "DELETE FROM building WHERE id = $1", testID)
}

func (s *RepositoriesSuite) TestBuildingRepository() {
	repo := NewBuilding(s.psql)
	buildings, err := repo.All(context.Background())
	s.NoError(err)
	s.NotEmpty(buildings)
}
