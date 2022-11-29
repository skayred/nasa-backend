package utils

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"panim.one/nasa/models"
)

var NasaTimeFormat = "2006-01-02"

type NASAServiceInterface interface {
	GetAsteroids(*logrus.Logger, time.Time, string, NASAAccessorInterface) (*models.NASAFeed, error)
}

type NASAService struct{}

func (NASAService) GetAsteroids(logger *logrus.Logger, at time.Time, apiKey string, accessor NASAAccessorInterface) (*models.NASAFeed, error) {
	// We take the asteroid on the daily basis
	to := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)
	to = to.AddDate(0, 0, 1)
	body, err := accessor.GetNASAJSON(logger, at, to, apiKey)

	if err != nil {
		logger.
			WithField("function", "GetAsteroids").
			WithField("step", "read").
			Error(err, "error during HTTP body read")

		return nil, err
	}

	feedObj := models.NASAFeed{}
	err = json.Unmarshal(body, &feedObj)

	if err != nil {
		logger.
			WithField("function", "GetAsteroids").
			WithField("step", "read").
			Error(err, "error during body JSON parse")

		return nil, err
	}

	return &feedObj, nil
}
