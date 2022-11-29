package utils

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var fetchNASAJSON func() ([]byte, error)

type testAccessor struct{}

func (acc testAccessor) GetNASAJSON(*logrus.Logger, time.Time, time.Time, string) ([]byte, error) {
	return fetchNASAJSON()
}

func TestNasaAPICall(t *testing.T) {
	logger, _ := test.NewNullLogger()

	fetchNASAJSON = func() ([]byte, error) { return []byte(properNASAJSON), nil }

	parser := NASAService{}
	asteroids, err := parser.GetAsteroids(logger, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "", testAccessor{})

	if err != nil {
		t.Fatal(err, "Parser shouldn't throw an error on the valid JSON!")
	}

	if len(asteroids.NearEarthObjects) != 1 {
		t.Fatalf("Expected 1 asteroid parsed, got %d", len(asteroids.NearEarthObjects))
	}

	asteroid := asteroids.NearEarthObjects["2015-09-08"][0]

	if *asteroid.ID != "2465633" ||
		*asteroid.Name != "465633 (2009 JR5)" ||
		asteroid.CloseApproachDates[0].CloseAt != 1441744080000 ||
		*asteroid.CloseApproachDates[0].MissDistance.KilometersAway != "45290438.204452618" {
		t.Fatal("Incorrect asteroid values parsed!")
	}
}

func TestNasaAPIErrorCall(t *testing.T) {
	logger, _ := test.NewNullLogger()

	fetchNASAJSON = func() ([]byte, error) { return []byte(""), nil }

	parser := NASAService{}
	_, err := parser.GetAsteroids(logger, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "", testAccessor{})

	if err == nil {
		t.Fatal(err, "Parser should throw an error on the empty JSON!")
	}
}
