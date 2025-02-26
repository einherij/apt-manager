package routes

import (
	"io"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/volatiletech/null/v8"

	"github.com/einherij/apt-manager/models"
)

func (s *RoutesSuite) TestBuildingHandler() {
	s.Run("All", func() {
		s.mockBuildingRepository.EXPECT().All(gomock.Any()).Return(models.BuildingSlice{
			{
				ID:      1,
				Name:    null.StringFrom("TestBuildingName"),
				Address: null.StringFrom("TestBuildingAddress"),
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/buildings", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)

		content, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.NoError(resp.Body.Close())
		s.JSONEq(`
			[{
				"id":1,
				"name":"TestBuildingName",
				"address":"TestBuildingAddress"
			}]`,
			string(content))
	})

	s.Run("Find", func() {
		s.mockBuildingRepository.EXPECT().Find(gomock.Any(), 1).Return(&models.Building{
			ID:      1,
			Name:    null.StringFrom("TestBuildingName"),
			Address: null.StringFrom("TestBuildingAddress"),
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/buildings/1", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)

		content, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.NoError(resp.Body.Close())
		s.JSONEq(`
			{
				"id":1,
				"name":"TestBuildingName",
				"address":"TestBuildingAddress"
			}`,
			string(content))
	})

	s.Run("Upsert", func() {
		s.mockBuildingRepository.EXPECT().Upsert(gomock.Any(), &models.Building{
			ID:      1,
			Name:    null.StringFrom("TestBuildingName"),
			Address: null.StringFrom("TestBuildingAddress"),
		}).Return(nil)
		body := `{"id":1,"name":"TestBuildingName","address":"TestBuildingAddress"}`
		req := httptest.NewRequest(fiber.MethodPost, "/buildings", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.NoError(err)
		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete", func() {
		s.mockBuildingRepository.EXPECT().Delete(gomock.Any(), 1).Return(nil)
		req := httptest.NewRequest(fiber.MethodDelete, "/buildings/1", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)
		s.Equal(fiber.StatusOK, resp.StatusCode)
	})
}
