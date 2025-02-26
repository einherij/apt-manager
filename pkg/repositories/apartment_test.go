package repositories

import (
	"context"

	"github.com/einherij/apt-manager/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
	s.NoError(building.Insert(context.Background(), s.psql, boil.Infer()))
	s.NoError(apartment.Insert(context.Background(), s.psql, boil.Infer()))

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
		for i := range apartments {
			s.Equal(testID, apartments[i].BuildingID.Int)
		}
	})

	s.Run("Delete", func() {
		s.NoError(repo.Delete(context.Background(), testID))
	})

	_, err := building.Delete(context.Background(), s.psql)
	s.NoError(err)
}
