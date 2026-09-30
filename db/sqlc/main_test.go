package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())

}
