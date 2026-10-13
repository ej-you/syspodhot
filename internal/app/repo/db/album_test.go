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

const _testDSN = "mongodb://127.0.0.1:27027/"
const _testDB = "syspodhot_test"

func TestAlbum_Create(t *testing.T) {
	repo, err := testAlbumRepo()
	require.NoError(t, err, "create album repo")

	album := &entity.Album{
		ID:         uuid.New(),
		Photo:      "file:///home/test/test.png",
		Title:      "Рич&гай",
		Type:       "альбом",
		Annotation: "Карты, деньги, два ствола",
		ReleasedAt: time.Now(),
		Streaming: entity.Streaming{
			VKLink:     "#",
			YandexLink: "#",
		},
		Songs: []entity.Song{
			{
				ID:       uuid.New(),
				Title:    "Карты",
				Duration: 224,
			},
			{
				ID:       uuid.New(),
				Title:    "Деньги",
				Duration: 296,
			},
			{
				ID:       uuid.New(),
				Title:    "Ствол",
				Duration: 189,
			},

			{
				ID:       uuid.New(),
				Title:    "Ствол 2",
				Duration: 177,
			},
		},
	}

	err = repo.Create(context.Background(), album)
	require.NoError(t, err, "create album record")
	t.Log("Created successfully")
}

func TestAlbum_GetAll(t *testing.T) {
	repo, err := testAlbumRepo()
	require.NoError(t, err, "create album repo")

	albums, err := repo.GetAll(context.Background())
	require.NoError(t, err, "get all albums")

	t.Log("Gotten albums:")
	for _, album := range albums {
		t.Logf("%+v", album)
	}
}

func TestAlbum_UpdateFull(t *testing.T) {
	repo, err := testAlbumRepo()
	require.NoError(t, err, "create album repo")

	album := &entity.Album{
		ID:         uuid.MustParse("32922206-3abd-48e7-8ab9-e08861ed7b61"),
		Photo:      "file:///home/test/test.png",
		Title:      "РичарТ",
		Type:       "альбом",
		Annotation: "Королевская классика",
		ReleasedAt: time.Now(),
		Streaming: entity.Streaming{
			VKLink:     "#",
			YandexLink: "#",
		},
		Songs: []entity.Song{
			{
				ID:       uuid.New(),
				Title:    "Зелено - зелено",
				Duration: 202,
			},
			{
				ID:       uuid.New(),
				Title:    "Тьма тьмущая",
				Duration: 252,
			},
		},
	}

	err = repo.UpdateFull(context.Background(), album)
	require.NoError(t, err, "update album")
}

func TestAlbum_Delete(t *testing.T) {
	repo, err := testAlbumRepo()
	require.NoError(t, err, "create album repo")

	id := uuid.MustParse("c82649cb-80b0-4605-9ea5-98c9a6689d98")
	err = repo.DeleteByID(context.Background(), id)
	require.NoError(t, err, "delete album")
}

func testAlbumRepo() (*AlbumRepoDB, error) {
	dbInst, err := db.New(_testDSN, _testDB)
	if err != nil {
		return nil, err
	}
	return NewAlbumRepoDB(dbInst), nil
}
