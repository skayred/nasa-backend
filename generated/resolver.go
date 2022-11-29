package generated

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"panim.one/nasa/models"
)

type Resolver struct{}

// // foo
func (r *proximityEventResolver) Asteroid(ctx context.Context, obj *models.ProximityEvent) (*models.Asteroid, error) {
	panic("not implemented")
}

// // foo
func (r *proximityEventResolver) MissDistance(ctx context.Context, obj *models.ProximityEvent) (string, error) {
	panic("not implemented")
}

// // foo
func (r *proximityEventResolver) HappenedAt(ctx context.Context, obj *models.ProximityEvent) (string, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) ClosestAsteroids(ctx context.Context, from string, to string, amount int) ([]*models.ProximityEvent, error) {
	panic("not implemented")
}

// ProximityEvent returns ProximityEventResolver implementation.
func (r *Resolver) ProximityEvent() ProximityEventResolver { return &proximityEventResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type proximityEventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
