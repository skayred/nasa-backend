package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type AsteroidServiceInterface interface {
	AsteroidByID(int) (*Asteroid, error)
	UpsertAsteroid(APIAsteroid) (*Asteroid, error)
	UpsertProximityEvent(*Asteroid, CloseApproachDate) (*ProximityEvent, error)
	ProximityEventsBetweenDates(time.Time, time.Time, int) ([]*ProximityEvent, error)
	WasDaySyncronized(time.Time) bool
}

type AsteroidService struct{}

func (AsteroidService) AsteroidByID(id int) (*Asteroid, error) {
	asteroid := &Asteroid{}

	err := PG.
		NewSelect().
		Model(asteroid).
		Where("id = ?", id).
		Scan(context.Background())

	return asteroid, err
}

func (AsteroidService) UpsertAsteroid(apiAsteroid APIAsteroid) (*Asteroid, error) {
	asteroid := &Asteroid{NasaID: *apiAsteroid.ID, Name: *apiAsteroid.Name}

	_, err := PG.
		NewInsert().
		Model(asteroid).
		On("CONFLICT (nasa_id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Returning("*").
		Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return asteroid, nil
}

func (AsteroidService) UpsertProximityEvent(asteroid *Asteroid, approachDate CloseApproachDate) (*ProximityEvent, error) {
	if asteroid == nil {
		panic("Asteroid cannot be nil when inserting the proximity event!")
	}

	if approachDate.MissDistance.KilometersAway == nil {
		panic(fmt.Sprintf("Asteroid miss distance is nil! Asteroid ID: %d, event date: %d", asteroid.ID, approachDate.CloseAt))
	}

	happenedAt := time.Unix(approachDate.CloseAt/1000, 0)
	missDistance, err := strconv.ParseFloat(*approachDate.MissDistance.KilometersAway, 64)

	if err != nil {
		panic(fmt.Sprintf("Couldn't parse the miss distance! Value was %s", *approachDate.MissDistance.KilometersAway))
	}

	existingEvents := []*ProximityEvent{}

	err = PG.
		NewSelect().
		Model(&existingEvents).
		Where(`asteroid_id = ? AND happened_at = ?`, asteroid.ID, happenedAt).
		Scan(context.Background())

	// Already exists - no need to save
	if len(existingEvents) > 0 {
		return existingEvents[0], nil
	}

	proximityEvent := &ProximityEvent{
		AsteroidID:   asteroid.ID,
		HappenedAt:   happenedAt,
		MissDistance: missDistance,
	}

	_, err = PG.
		NewInsert().
		Model(proximityEvent).
		Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return proximityEvent, nil
}

func (AsteroidService) ProximityEventsBetweenDates(from time.Time, to time.Time, amount int) ([]*ProximityEvent, error) {
	proximityEvents := []*ProximityEvent{}

	err := PG.
		NewSelect().
		Model(&proximityEvents).
		Where(`happened_at >= ? AND happened_at <= ?`, from, to).
		Limit(amount).
		OrderExpr("miss_distance ASC").
		Scan(context.Background())

	return proximityEvents, err
}

func (AsteroidService) WasDaySyncronized(at time.Time) bool {
	count, err := PG.
		NewSelect().
		Model((*ProximityEvent)(nil)).
		Where(`happened_at::date = ?::date`, time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)).
		Count(context.Background())

	if err == sql.ErrNoRows || (count == 0) {
		return false
	}

	if err != nil {
		panic("Error when checking for proximity events day existence!")
	}

	return true
}
