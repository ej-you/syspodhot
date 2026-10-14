package db

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"syspodhot/internal/app/entity"
	"syspodhot/internal/pkg/db"
)

func TestConcert_Create(t *testing.T) {
	repo, err := testConcertRepo()
	require.NoError(t, err, "create concert repo")

	concert := &entity.Concert{
		ID:         uuid.New(),
		Title:      "Новогодний концерт",
		Place:      "Нижний Новгород",
		Datetime:   time.Now().AddDate(0, 2, 12),
		TicketLink: "#",
		Photos: []string{
			"file:///home/test/concert1_1.png",
			"file:///home/test/concert1_2.png",
		},
		Video: "#",
	}

	err = repo.Create(context.Background(), concert)
	require.NoError(t, err, "create concert record")
	t.Log("Created successfully")
}

func TestConcert_GetAll(t *testing.T) {
	repo, err := testConcertRepo()
	require.NoError(t, err, "create concert repo")

	concerts, err := repo.GetMany(context.Background(), false, false)
	require.NoError(t, err, "get all concerts")

	t.Log("Gotten concerts:")
	for _, concert := range concerts {
		t.Logf("%+v", concert)
	}
}

func testConcertRepo() (*ConcertRepoDB, error) {
	dbInst, err := db.New(_testDSN, _testDB)
	if err != nil {
		return nil, err
	}
	return NewConcertRepoDB(dbInst), nil
}
