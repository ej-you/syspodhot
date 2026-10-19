package entity

import (
	"fmt"
	"time"
	"uuid"

	"syspodhot/internal/pkg/utils/timefmt"
)

const (
	_minute     = 60
	_halfMinute = 30
)

// Album represents a music album with a slice of songs.
type Album struct {
	ID         uuid.UUID `json:"id"`
	Photo      string    `json:"photo"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	Annotation string    `json:"annotation"`
	ReleasedAt time.Time `json:"released_at"`
	Streaming  Streaming `json:"streaming"`
	Songs      []Song    `json:"songs"`
}

// TotalDuration returns human-readable duration of all songs like "40 мин".
func (a *Album) TotalDuration() string {
	var total int
	for _, song := range a.Songs {
		total += song.Duration
	}

	totalMins := total / _minute
	if total%_minute > _halfMinute {
		totalMins++
	}
	return fmt.Sprintf("%d мин", totalMins)
}

// ReleasedAtH returns human-readable release date like "04 июля 2026".
func (a *Album) ReleasedAtH() string {
	return timefmt.FmtRuDateFull(a.ReleasedAt)
}

// Streaming represents a streaming platform links for album.
type Streaming struct {
	VKLink     string `json:"vk"`
	YandexLink string `json:"yandex"`
}

// Song represents a one song.
type Song struct {
	Title string `json:"title"`
	// song duration in seconds
	Duration int `json:"duration"`
}

// DurationH returns human-readable duration of song like "HH:mm:ss" or "mm:ss".
func (s *Song) DurationH() string {
	return timefmt.FmtDuration(time.Duration(s.Duration) * time.Second)
}
