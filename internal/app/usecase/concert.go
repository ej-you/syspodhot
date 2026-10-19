package usecase

import (
	"context"
	"slices"
	"time"
	"uuid"

	"syspodhot/internal/app/entity"
	repodb "syspodhot/internal/app/repo/db"
)

// ConcertUC represents a concert usecase.
type ConcertUC struct {
	concertRepoDB *repodb.ConcertRepoDB
	locale        *time.Location
}

// NewConcertUC returns a new instance of ConcertUC.
func NewConcertUC(concertRepoDB *repodb.ConcertRepoDB, locale *time.Location) *ConcertUC {
	return &ConcertUC{
		concertRepoDB: concertRepoDB,
		locale:        locale,
	}
}

func (u *ConcertUC) Create(concert *entity.Concert) error {
	concert.ID = uuid.New()
	concert.Datetime = concert.Datetime.UTC()
	if err := u.concertRepoDB.Create(context.Background(), concert); err != nil {
		return err
	}
	concert.Datetime = concert.Datetime.In(u.locale)
	return nil
}

func (u *ConcertUC) Get(id uuid.UUID) (*entity.Concert, error) {
	concert, err := u.concertRepoDB.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	concert.Datetime = concert.Datetime.In(u.locale)
	return concert, nil
}

// GetComing returns all coming concerts.
func (u *ConcertUC) GetComing() ([]*entity.Concert, error) {
	concerts, err := u.concertRepoDB.GetMany(context.Background(), true, false)
	if err != nil {
		return nil, err
	}
	for idx := range concerts {
		concerts[idx].Datetime = concerts[idx].Datetime.In(u.locale)
	}
	return concerts, nil
}

// GetPastGroupedByYear returns all past concerts grouped by year.
func (u *ConcertUC) GetPastGroupedByYear() ([]entity.ConcertsByYear, error) {
	concerts, err := u.concertRepoDB.GetMany(context.Background(), false, true)
	if err != nil {
		return nil, err
	}

	// group concerts by year
	byYearMap := make(map[int][]*entity.Concert)
	var year int
	for _, concert := range concerts {
		// set timezone for concert datetime
		concert.Datetime = concert.Datetime.In(u.locale)
		year = concert.Year()
		byYearMap[year] = append(byYearMap[year], concert)
	}

	// collect slice of ConcertsByYear structs
	groups := make([]entity.ConcertsByYear, 0, len(byYearMap))
	for year, concerts := range byYearMap {
		groups = append(groups, entity.ConcertsByYear{Year: year, Concerts: concerts})
	}
	// sort by year
	slices.SortFunc(groups, func(a, b entity.ConcertsByYear) int {
		return b.Year - a.Year // DESC
	})

	return groups, nil
}

func (u *ConcertUC) Update(concert *entity.Concert) error {
	concert.Datetime = concert.Datetime.UTC()
	if err := u.concertRepoDB.UpdateFull(context.Background(), concert); err != nil {
		return err
	}
	concert.Datetime = concert.Datetime.In(u.locale)
	return nil
}

func (u *ConcertUC) Delete(id uuid.UUID) error {
	return u.concertRepoDB.DeleteByID(context.Background(), id) // err OR nil
}
