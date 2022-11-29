package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProximityEvent struct {
	bun.BaseModel `bun:"proximity_events"`

	ID           int       `json:"id" bun:"id,pk,autoincrement"`
	AsteroidID   int       `json:"asteroidID" bun:"asteroid_id"`
	Asteroid     *Asteroid `json:"asteroid" bun:"rel:belongs-to,join:asteroid_id=id"`
	MissDistance float64   `json:"missDistance" bun:"miss_distance"`
	HappenedAt   time.Time `json:"happenedAt" bun:"happened_at"`
}
