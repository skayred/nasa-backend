package main

import (
	"context"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"panim.one/nasa/generated"
	"panim.one/nasa/models"
	"panim.one/nasa/resolvers"
	"panim.one/nasa/utils"
)

var logger = logrus.New()

func main() {
	// Init singletone PG connection
	models.Init()
	router := chi.NewRouter()

	apiKey, ok := os.LookupEnv("NASA_API_KEY")

	if !ok {
		panic("NASA_API_KEY environment variable was not defined!")
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://panim.one", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	resolverImplementation := generated.Config{
		Resolvers: &resolvers.Resolver{
			AsteroidService: models.AsteroidService{},
			NASAService:     utils.NASAService{},
			NASAAccessor:    utils.NASAAccessor{},
			NASAAPIKey:      apiKey,
			Logger:          logger,
		},
	}

	srv := handler.New(generated.NewExecutableSchema(resolverImplementation))
	srv.AddTransport(transport.POST{})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		rc.DisableIntrospection = false

		logCtx := context.WithValue(ctx, "LOGGER", logger)
		return next(logCtx)
	})

	router.Handle("/graphql", srv)

	logger.Println("Server running at *:4000")
	logger.Fatal(http.ListenAndServe(":4000", router))
}
