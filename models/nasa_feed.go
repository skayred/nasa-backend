package models

type MissDistance struct {
	KilometersAway *string `json:"kilometers"`
}

type CloseApproachDate struct {
	CloseAt      int64        `json:"epoch_date_close_approach"`
	OrbitingBody *string      `json:"orbiting_body"`
	MissDistance MissDistance `json:"miss_distance"`
}

type APIAsteroid struct {
	ID                 *string             `json:"id"`
	Name               *string             `json:"name"`
	CloseApproachDates []CloseApproachDate `json:"close_approach_data"`
}

type NASAFeed struct {
	NearEarthObjects map[string][]APIAsteroid `json:"near_earth_objects"`
}
