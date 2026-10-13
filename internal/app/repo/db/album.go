package db

import (
	"context"
	goerrors "errors"
	"fmt"
	"uuid"

	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/pkg/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _collectionAlbum = (album{}).collectionName()

var _idField = "_id" // id field name in mongodb

// AlbumRepoDB represents an album DB repo.
type AlbumRepoDB struct {
	dbStorage *db.DB
}

// NewAlbumRepoDB returns a new instance of AlbumRepoDB.
func NewAlbumRepoDB(dbStorage *db.DB) *AlbumRepoDB {
	return &AlbumRepoDB{
		dbStorage: dbStorage,
	}
}

// Create creates a new album record.
func (r *AlbumRepoDB) Create(ctx context.Context, album *entity.Album) error {
	albumDTO := newAlbum(album)

	_, err := r.dbStorage.
		Collection(_collectionAlbum).
		InsertOne(ctx, albumDTO)
	if err != nil {
		return err
	}
	return nil
}

// GetByID returns album with given ID.
func (r *AlbumRepoDB) GetByID(ctx context.Context, id uuid.UUID) (*entity.Album, error) {
	album := &album{}

	err := r.dbStorage.
		Collection(_collectionAlbum).
		FindOne(ctx, bson.M{_idField: id}).
		Decode(album)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return album.toEntity(), nil
}

// GetAll returns all albums.
func (r *AlbumRepoDB) GetAll(ctx context.Context) ([]*entity.Album, error) {
	albumListRaw := make([]album, 0)

	opts := options.Find().SetSort(bson.D{{Key: "released_at", Value: -1}})
	cursor, err := r.dbStorage.
		Collection(_collectionAlbum).
		Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("get cursor: %w", err)
	}
	if err := cursor.All(ctx, &albumListRaw); err != nil {
		return nil, fmt.Errorf("decode results: %w", err)
	}

	albumList := make([]*entity.Album, len(albumListRaw))
	for idx := range albumListRaw {
		albumList[idx] = albumListRaw[idx].toEntity()
	}
	return albumList, nil
}

// UpdateFull updates full album data by its ID.
func (r *AlbumRepoDB) UpdateFull(ctx context.Context, album *entity.Album) error {
	albumDTO := newAlbum(album)

	result, err := r.dbStorage.
		Collection(_collectionAlbum).
		ReplaceOne(ctx, bson.M{_idField: albumDTO.ID}, albumDTO)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// DeleteByID deletes album with given ID.
func (r *AlbumRepoDB) DeleteByID(ctx context.Context, id uuid.UUID) error {
	album := &album{}

	err := r.dbStorage.
		Collection(_collectionAlbum).
		FindOneAndDelete(ctx, bson.M{_idField: id}).
		Decode(album)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrNotFound
	}
	return err // err OR nil
}
