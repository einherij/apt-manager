package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	mock_routes "github.com/einherij/apt-manager/pkg/mocks/routes"
)

type RoutesSuite struct {
	suite.Suite

	ctrl                    *gomock.Controller
	mockBuildingRepository  *mock_routes.MockBuildingRepository
	mockApartmentRepository *mock_routes.MockApartmentRepository

	app              *fiber.App
	buildingHandler  *BuildingHandler
	apartmentHandler *ApartmentHandler
}

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RoutesSuite))
}

func (s *RoutesSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockBuildingRepository = mock_routes.NewMockBuildingRepository(s.ctrl)
	s.mockApartmentRepository = mock_routes.NewMockApartmentRepository(s.ctrl)

	s.app = fiber.New()
	s.buildingHandler = NewBuildingHandler(s.mockBuildingRepository)
	s.apartmentHandler = NewApartmentHandler(s.mockApartmentRepository)
	RegisterRoutes(s.app, s.buildingHandler, s.apartmentHandler)
}

func (s *RoutesSuite) TearDownTest() {
	s.ctrl.Finish()
}
