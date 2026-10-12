package entity

import (
	"time"
	"uuid"

	"syspodhot/internal/pkg/utils/timefmt"
)

// Album represents a music album with a slice of songs.
type Album struct {
	ID         uuid.UUID `json:"id"`
	Photo      string    `json:"photo"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	Annotation string    `json:"annotation"`
	ReleasedAt time.Time `json:"released_at"`
	Songs      []Song    `json:"songs"`
}

// TotalDuration returns human-readable duration of all songs like "HH:mm:ss" or "mm:ss".
func (a *Album) TotalDuration() string {
	var total int
	for _, song := range a.Songs {
		total += song.Duration
	}
	return timefmt.FmtDuration(time.Duration(total) * time.Second)
}

// ReleasedAtH returns human-readable release date like "04 июля 2026".
func (a *Album) ReleasedAtH() string {
	return timefmt.FmtRuDateFull(a.ReleasedAt)
}

// Song represents a one song.
type Song struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	// song duration in seconds
	Duration int `json:"duration"`
}

// DurationH returns human-readable duration of song like "HH:mm:ss" or "mm:ss".
func (s *Song) DurationH() string {
	return timefmt.FmtDuration(time.Duration(s.Duration) * time.Second)
}
