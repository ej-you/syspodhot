package usecase

import (
	"context"
	"syspodhot/internal/app/entity"
	repodb "syspodhot/internal/app/repo/db"
	"uuid"
)

// ConcertUC represents a concert usecase.
type ConcertUC struct {
	concertRepoDB *repodb.ConcertRepoDB
}

// NewConcertUC returns a new instance of ConcertUC.
func NewConcertUC(concertRepoDB *repodb.ConcertRepoDB) *ConcertUC {
	return &ConcertUC{
		concertRepoDB: concertRepoDB,
	}
}

func (u *ConcertUC) Create(concert *entity.Concert) error {
	concert.ID = uuid.New()
	concert.Datetime = concert.Datetime.UTC()
	return u.concertRepoDB.Create(context.Background(), concert) // err OR nil
}

func (u *ConcertUC) Get(id uuid.UUID) (*entity.Concert, error) {
	concert, err := u.concertRepoDB.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return concert, nil
}

// GetComming returns all comming concerts.
func (u *ConcertUC) GetComming() ([]*entity.Concert, error) {
	concerts, err := u.concertRepoDB.GetMany(context.Background(), true, false)
	if err != nil {
		return nil, err
	}
	return concerts, nil
}

// GetPast returns all past concerts.
func (u *ConcertUC) GetPast() ([]*entity.Concert, error) {
	concerts, err := u.concertRepoDB.GetMany(context.Background(), false, true)
	if err != nil {
		return nil, err
	}
	return concerts, nil
}

func (u *ConcertUC) Update(concert *entity.Concert) error {
	concert.Datetime = concert.Datetime.UTC()
	return u.concertRepoDB.UpdateFull(context.Background(), concert) // err OR nil
}

func (u *ConcertUC) Delete(id uuid.UUID) error {
	return u.concertRepoDB.DeleteByID(context.Background(), id) // err OR nil
}
