package models

import (
	"fmt"
	"os"
)

var dbHostEnv = "DB_HOST"
var dbNameEnv = "DB_NAME"
var dbUsernameEnv = "DB_USERNAME"
var dbPasswordEnv = "DB_PASSWORD"

var requiredEnvs = []string{dbHostEnv, dbNameEnv, dbUsernameEnv, dbPasswordEnv}

func BuildPGConnectionString() string {
	for _, envName := range requiredEnvs {
		_, ok := os.LookupEnv(envName)

		if !ok {
			panic(fmt.Sprintf("Missing environment variable: %s!", envName))
		}
	}

	dbHost := os.Getenv(dbHostEnv)
	dbName := os.Getenv(dbNameEnv)
	dbUsername := os.Getenv(dbUsernameEnv)
	dbPassword := os.Getenv(dbPasswordEnv)

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUsername, dbPassword, dbHost, dbName)
}
