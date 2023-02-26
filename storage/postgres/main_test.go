package postgres

import (
	"context"
	"fmt"
	"go_auth_api_gateway/config"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/manveru/faker"
	"github.com/stretchr/testify/assert"
)

var (
	db       *pgxpool.Pool
	fakeData *faker.Faker
)

func createRandomId(t *testing.T) string {
	id, err := uuid.NewRandom()
	assert.NoError(t, err)
	return id.String()
}

func TestMain(m *testing.M) {
	cfg := config.Load()
	// conf, err := pgxpool.ParseConfig(fmt.Sprintf(
	//   "postgres://%s:%s@%s:%d/%s?sslmode=disable",
	//   cfg.PostgresUser,
	//   cfg.PostgresPassword,
	//   cfg.PostgresHost,
	//   cfg.PostgresPort,
	//   cfg.PostgresDatabase,
	// ))
	// if err != nil {
	//   panic(err)
	// }

	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		panic(err)
	}

	fakeData, _ = faker.New("en")
	conf.MaxConns = cfg.PostgresMaxConnections

	dbPool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		panic(err)
	}

	db = dbPool

	os.Exit(m.Run())
}
