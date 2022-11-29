package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type NASAAccessorInterface interface {
	GetNASAJSON(*logrus.Logger, time.Time, time.Time, string) ([]byte, error)
}

type NASAAccessor struct{}

func (acc NASAAccessor) GetNASAJSON(logger *logrus.Logger, from time.Time, to time.Time, apiKey string) ([]byte, error) {
	url := fmt.Sprintf("https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s", from.Format(NasaTimeFormat), to.Format(NasaTimeFormat), apiKey)
	resp, err := http.Get(url)

	if err != nil {
		logger.
			WithField("function", "GetNASAJSON").
			WithField("step", "http").
			Error(err, "error during HTTP call")

		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
