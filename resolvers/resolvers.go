package resolvers

import (
	"github.com/sirupsen/logrus"
	"panim.one/nasa/generated"
	"panim.one/nasa/models"
	"panim.one/nasa/utils"
)

type Resolver struct {
	AsteroidService models.AsteroidServiceInterface
	NASAService     utils.NASAServiceInterface
	NASAAccessor    utils.NASAAccessorInterface
	Logger          *logrus.Logger
	NASAAPIKey      string
}

// ProximityEvent returns ProximityEventResolver implementation.
func (r *Resolver) ProximityEvent() generated.ProximityEventResolver {
	return &proximityEventResolver{r}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type proximityEventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
