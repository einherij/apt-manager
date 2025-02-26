package routes

import (
	"io"
	"net/http/httptest"
	"strings"

	"github.com/einherij/apt-manager/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/volatiletech/null/v8"
)

func (s *RoutesSuite) TestApartmentHandler() {
	s.Run("All", func() {
		s.mockApartmentRepository.EXPECT().All(gomock.Any()).Return(models.ApartmentSlice{
			{
				ID:         1,
				Number:     null.StringFrom("10"),
				BuildingID: null.IntFrom(1),
				Floor:      null.IntFrom(4),
				SQMeters:   null.IntFrom(30),
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/apartments", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)

		content, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.NoError(resp.Body.Close())
		s.JSONEq(`
			[{
				"id":1,
				"building_id":1,
				"number":"10",
				"floor":4,
				"sq_meters":30
			}]`,
			string(content))
	})

	s.Run("Find", func() {
		s.mockApartmentRepository.EXPECT().Find(gomock.Any(), 1).Return(&models.Apartment{
			ID:         1,
			Number:     null.StringFrom("10"),
			BuildingID: null.IntFrom(1),
			Floor:      null.IntFrom(4),
			SQMeters:   null.IntFrom(30),
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/apartments/1", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)

		content, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.NoError(resp.Body.Close())
		s.JSONEq(`
			{
				"id":1,
				"building_id":1,
				"number":"10",
				"floor":4,
				"sq_meters":30
			}`,
			string(content))
	})

	s.Run("FindByBuildingID", func() {
		s.mockApartmentRepository.EXPECT().FindByBuildingID(gomock.Any(), 1).Return(models.ApartmentSlice{
			{
				ID:         1,
				Number:     null.StringFrom("10"),
				BuildingID: null.IntFrom(1),
				Floor:      null.IntFrom(4),
				SQMeters:   null.IntFrom(30),
			},
		}, nil)
		req := httptest.NewRequest(fiber.MethodGet, "/apartments/building/1", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)

		content, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.NoError(resp.Body.Close())
		s.JSONEq(`
			[{
				"id":1,
				"building_id":1,
				"number":"10",
				"floor":4,
				"sq_meters":30
			}]`,
			string(content))
	})

	s.Run("Upsert", func() {
		s.mockApartmentRepository.EXPECT().Upsert(gomock.Any(), &models.Apartment{
			ID:         1,
			Number:     null.StringFrom("10"),
			BuildingID: null.IntFrom(1),
			Floor:      null.IntFrom(4),
			SQMeters:   null.IntFrom(30),
		}).Return(nil)
		body := `{"id":1,"number":"10","building_id":1,"floor":4,"sq_meters":30}`
		req := httptest.NewRequest(fiber.MethodPost, "/apartments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		s.NoError(err)
		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete", func() {
		s.mockApartmentRepository.EXPECT().Delete(gomock.Any(), 1).Return(nil)
		req := httptest.NewRequest(fiber.MethodDelete, "/apartments/1", nil)

		resp, err := s.app.Test(req)
		s.NoError(err)
		s.Equal(fiber.StatusOK, resp.StatusCode)
	})
}
