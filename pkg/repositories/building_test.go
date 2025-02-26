package repositories

import (
	"context"

	"github.com/einherij/apt-manager/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *RepositoriesSuite) TestBuildingRepository() {
	repo := NewBuilding(s.psql)
	const testID = 1

	building := &models.Building{
		ID:      testID,
		Name:    null.StringFrom("TestBuildingName"),
		Address: null.StringFrom("TestBuildingAddress"),
	}
	s.NoError(building.Insert(context.Background(), s.psql, boil.Infer()))

	s.Run("Upsert", func() {
		building.Name = null.StringFrom("NewBuildingName")
		err := repo.Upsert(context.Background(), building)
		s.NoError(err)
		b, err := repo.Find(context.Background(), testID)
		s.NoError(err)
		s.Equal(building.Name, b.Name)
	})

	s.Run("All", func() {
		buildings, err := repo.All(context.Background())
		s.NoError(err)
		s.NotEmpty(buildings)
	})

	s.Run("Find", func() {
		b, err := repo.Find(context.Background(), testID)
		s.NoError(err)
		s.NotNil(b)
	})

	s.Run("Delete", func() {
		s.NoError(repo.Delete(context.Background(), testID))
	})
}
