package models

import "github.com/uptrace/bun"

type Asteroid struct {
	bun.BaseModel `bun:"asteroids"`

	ID              int              `json:"id" bun:"id,pk,autoincrement"`
	NasaID          string           `json:"nasaID" bun:"nasa_id,unique"`
	Name            string           `json:"name" bun:"name"`
	ProximityEvents []ProximityEvent `json:"proximityEvents" bun:"rel:has-many,join:id=asteroid_id"`
}
