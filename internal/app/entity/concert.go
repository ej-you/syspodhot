package entity

import (
	"time"
	"uuid"

	"syspodhot/internal/pkg/utils/timefmt"
)

const _concertTimeFmt = "15:04" // HH:mm

// Concert represents a coming/past music concert.
type Concert struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Place      string    `json:"place"`
	Datetime   time.Time `json:"datetime"`
	TicketLink string    `json:"ticket_link"`

	Photos []string `json:"photos"`
	Video  string   `json:"video"`
}

// DatetimeH returns human-readable concert date like "04 июля 2026".
func (c *Concert) DatetimeH() string {
	return timefmt.FmtRuDateFull(c.Datetime)
}

// DatetimeH returns concert month like "января".
func (c *Concert) Month() string {
	return timefmt.FmtRuMonth(c.Datetime)
}

// DatetimeH returns concert day of month.
func (c *Concert) Day() int {
	return c.Datetime.Day()
}

// DatetimeH returns human-readable concert weekday like "суббота".
func (c *Concert) WeekDay() string {
	return timefmt.FmtRuWeekDay(c.Datetime)
}

// DatetimeH returns human-readable concert time like "HH:mm".
func (c *Concert) Time() string {
	return c.Datetime.Format(_concertTimeFmt)
}
