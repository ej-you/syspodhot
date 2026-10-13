package db

import (
	"context"
	goerrors "errors"
	"fmt"
	"syspodhot/internal/app/entity"
	"syspodhot/internal/app/errors"
	"syspodhot/internal/pkg/db"
	"time"
	"uuid"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _collectionConcert = (concert{}).collectionName()

// ConcertRepoDB represents a concert DB repo.
type ConcertRepoDB struct {
	dbStorage *db.DB
}

// NewConcertRepoDB returns a new instance of ConcertRepoDB.
func NewConcertRepoDB(dbStorage *db.DB) *ConcertRepoDB {
	return &ConcertRepoDB{
		dbStorage: dbStorage,
	}
}

// Create creates a new concert record.
func (r *ConcertRepoDB) Create(ctx context.Context, concert *entity.Concert) error {
	concertDTO := newConcert(concert)

	_, err := r.dbStorage.
		Collection(_collectionConcert).
		InsertOne(ctx, concertDTO)
	if err != nil {
		return err
	}
	return nil
}

// GetByID returns concert with given ID.
func (r *ConcertRepoDB) GetByID(ctx context.Context, id uuid.UUID) (*entity.Concert, error) {
	concert := &concert{}

	err := r.dbStorage.
		Collection(_collectionConcert).
		FindOne(ctx, bson.M{_idField: id}).
		Decode(concert)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return concert.toEntity(), nil
}

// GetAll returns all concerts. If commingOnly is true it returns only comming concerts.
func (r *ConcertRepoDB) GetMany(ctx context.Context,
	commingOnly, pastOnly bool) ([]*entity.Concert, error) {

	concertListRaw := make([]concert, 0)

	// apply filter if commingOnly is true
	filter := bson.M{}
	if commingOnly && !pastOnly {
		filter = bson.M{
			"datetime": bson.M{"$gt": time.Now().UTC()},
		}
	} else if pastOnly && !commingOnly {
		filter = bson.M{
			"datetime": bson.M{"$lte": time.Now().UTC()},
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "datetime", Value: -1}})
	cursor, err := r.dbStorage.
		Collection(_collectionConcert).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("get cursor: %w", err)
	}
	if err := cursor.All(ctx, &concertListRaw); err != nil {
		return nil, fmt.Errorf("decode results: %w", err)
	}

	concertList := make([]*entity.Concert, len(concertListRaw))
	for idx := range concertListRaw {
		concertList[idx] = concertListRaw[idx].toEntity()
	}
	return concertList, nil
}

// UpdateFull updates full concert data by its ID.
func (r *ConcertRepoDB) UpdateFull(ctx context.Context, concert *entity.Concert) error {
	concertDTO := newConcert(concert)

	result, err := r.dbStorage.
		Collection(_collectionConcert).
		ReplaceOne(ctx, bson.M{_idField: concertDTO.ID}, concertDTO)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// DeleteByID deletes concert with given ID.
func (r *ConcertRepoDB) DeleteByID(ctx context.Context, id uuid.UUID) error {
	concert := &concert{}

	err := r.dbStorage.
		Collection(_collectionConcert).
		FindOneAndDelete(ctx, bson.M{_idField: id}).
		Decode(concert)
	if goerrors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrNotFound
	}
	return err // err OR nil
}
