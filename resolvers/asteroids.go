package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"panim.one/nasa/models"
	"panim.one/nasa/utils"
)

func failOnError(logger *logrus.Logger, err error, place string) {
	if err != nil {
		logger.Error(err, "Error happened, defaulting to panic!", place)
		panic(err)
	}
}

func (resolver *proximityEventResolver) Asteroid(ctx context.Context, obj *models.ProximityEvent) (*models.Asteroid, error) {
	if obj == nil {
		return nil, fmt.Errorf("Object was null!")
	}

	return resolver.AsteroidService.AsteroidByID(obj.AsteroidID)
}

func (resolver *proximityEventResolver) MissDistance(ctx context.Context, obj *models.ProximityEvent) (string, error) {
	if obj == nil {
		return "0", fmt.Errorf("Object was null!")
	}

	return fmt.Sprintf("%f", obj.MissDistance), nil
}

func (resolver *proximityEventResolver) HappenedAt(ctx context.Context, obj *models.ProximityEvent) (string, error) {
	if obj == nil {
		return time.Now().Format(utils.NasaTimeFormat), fmt.Errorf("Object was null!")
	}

	return obj.HappenedAt.Format(utils.NasaTimeFormat), nil
}

func (resolver *queryResolver) ClosestAsteroids(ctx context.Context, fromString string, toString string, amount int) ([]*models.ProximityEvent, error) {
	resolver.Logger.
		WithField("from", fromString).
		WithField("to", toString).
		Info(fmt.Sprintf("Incodming request for the closest asteroids between %s and %s", fromString, toString))

	now := time.Now().UTC()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from, fromErr := time.Parse(utils.NasaTimeFormat, fromString)
	to, toErr := time.Parse(utils.NasaTimeFormat, toString)

	if fromErr != nil || toErr != nil {
		panic("Failed to parse input date! Format expected is: yyyy-MM-dd")
	}

	if from.After(now) {
		from = startOfToday
	}

	if to.After(now) {
		to = startOfToday
	}

	daysAffected := utils.DaysBetweenDates(from, to)

	resolver.Logger.
		WithField("from", fromString).
		WithField("to", toString).
		WithField("days_amount", len(daysAffected)).
		Info(fmt.Sprintf("Days affected by the request: %d", len(daysAffected)))

	for _, day := range daysAffected {
		if !resolver.AsteroidService.WasDaySyncronized(day) {
			resolver.Logger.
				WithField("from", fromString).
				WithField("to", toString).
				WithField("day", day).
				Info(fmt.Sprintf("Day %s was not synced", day))

			feed, err := resolver.NASAService.GetAsteroids(resolver.Logger, day, resolver.NASAAPIKey, resolver.NASAAccessor)
			failOnError(resolver.Logger, err, "asteroids retrieval")

			for _, asteroids := range feed.NearEarthObjects {
				for _, apiAsteroid := range asteroids {
					resolver.Logger.
						WithField("from", fromString).
						WithField("to", toString).
						WithField("day", day).
						WithField("asteroid_id", *apiAsteroid.ID).
						WithField("asteroid_name", *apiAsteroid.Name).
						Info(fmt.Sprintf("Asteroid downloaded: %s, name is %s", *apiAsteroid.ID, *apiAsteroid.Name))

					asteroid, err := resolver.AsteroidService.UpsertAsteroid(apiAsteroid)
					failOnError(resolver.Logger, err, "asteroid upsertion")

					for _, event := range apiAsteroid.CloseApproachDates {
						resolver.Logger.
							WithField("from", fromString).
							WithField("to", toString).
							WithField("day", day).
							WithField("asteroid_id", *apiAsteroid.ID).
							WithField("asteroid_name", *apiAsteroid.Name).
							Info(fmt.Sprintf("Proximity event downloaded: distance is %s, epoch time is %d", *event.MissDistance.KilometersAway, event.CloseAt))

						_, err := resolver.AsteroidService.UpsertProximityEvent(asteroid, event)
						failOnError(resolver.Logger, err, "proximity event upsertion")
					}
				}
			}
		} else {
			resolver.Logger.
				WithField("from", fromString).
				WithField("to", toString).
				WithField("day", day).
				Info(fmt.Sprintf("Day %s was already synced", day))
		}
	}

	return resolver.AsteroidService.ProximityEventsBetweenDates(from, to, amount)
}
