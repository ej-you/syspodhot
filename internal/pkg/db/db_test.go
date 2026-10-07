package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnDB(t *testing.T) {
	dsn := "mongodb://127.0.0.1:27017/"
	dbName := "syspodhot_test"

	db, err := New(dsn, dbName)
	require.NoError(t, err, "Connect to mongodb")
	defer db.Close()

	t.Logf("Conn to db with name: %v", db.Name())
}
