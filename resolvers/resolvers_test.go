package resolvers

import (
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"panim.one/nasa/generated"
	"panim.one/nasa/models"
	"panim.one/nasa/utils"
)

type mockAsteroidService struct {
	UpsertAsteroidCalled  int
	UpsertProximityCalled int
}

type mockNASAService struct{}

type mockNASAAccessor struct{}

// Asteroid service mock
func (mockAsteroidService) AsteroidByID(int) (*models.Asteroid, error) {
	return &models.Asteroid{ID: 1}, nil
}

func (service mockAsteroidService) UpsertAsteroid(apiAsteroid models.APIAsteroid) (*models.Asteroid, error) {
	service.UpsertAsteroidCalled = service.UpsertAsteroidCalled + 1
	return &models.Asteroid{ID: 1}, nil
}

func (service mockAsteroidService) UpsertProximityEvent(asteroid *models.Asteroid, approachDate models.CloseApproachDate) (*models.ProximityEvent, error) {
	service.UpsertProximityCalled = service.UpsertProximityCalled + 1
	return nil, nil

}
func (mockAsteroidService) ProximityEventsBetweenDates(from time.Time, to time.Time, amount int) ([]*models.ProximityEvent, error) {
	return []*models.ProximityEvent{
		&models.ProximityEvent{
			ID:           1,
			MissDistance: 123,
			HappenedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}, nil
}

func (mockAsteroidService) WasDaySyncronized(at time.Time) bool {
	return true
}

// NASA Service mock
func (mockNASAService) GetAsteroids(*logrus.Logger, time.Time, string, utils.NASAAccessorInterface) (*models.NASAFeed, error) {
	return nil, nil
}

// NASA accessor
func (mockNASAAccessor) GetNASAJSON(*logrus.Logger, time.Time, time.Time, string) ([]byte, error) {
	return []byte{}, nil
}

type TestAsteroid struct {
	ID string
}

type TestProximityEvent struct {
	ID           string
	MissDistance string
	HappenedAt   string
	Asteroid     TestAsteroid
}

type Response struct {
	ClosestAsteroids []TestProximityEvent
}

func TestMutationResolver_ValidateAccessToken(t *testing.T) {
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		serviceMock := mockAsteroidService{UpsertAsteroidCalled: 0, UpsertProximityCalled: 0}
		nasaMock := mockNASAService{}
		accessorMock := mockNASAAccessor{}

		testResolvers := Resolver{
			AsteroidService: serviceMock,
			NASAService:     nasaMock,
			NASAAccessor:    accessorMock,
			Logger:          logger,
			NASAAPIKey:      "",
		}

		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &testResolvers})))

		resp := Response{}
		query := `
		{
			closestAsteroids(from: "2020-01-01", to: "2020-01-07", amount: 3) {
			  id
			  missDistance
			  happenedAt
			  asteroid {
				id
			  }
			}
		  }
		`
		c.MustPost(query, &resp)
		require.Equal(t, "1", resp.ClosestAsteroids[0].ID)
		require.Equal(t, "123.000000", resp.ClosestAsteroids[0].MissDistance)
		require.Equal(t, "2020-01-01", resp.ClosestAsteroids[0].HappenedAt)
		require.Equal(t, "1", resp.ClosestAsteroids[0].Asteroid.ID)

		// As day is already fetched - no upserts should happen
		require.Equal(t, 0, serviceMock.UpsertAsteroidCalled)
		require.Equal(t, 0, serviceMock.UpsertProximityCalled)
	})

}
