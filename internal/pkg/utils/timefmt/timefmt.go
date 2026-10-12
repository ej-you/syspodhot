package timefmt

import (
	"fmt"
	"time"
)

var (
	_monthsRu = []string{
		"", "января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
	_weekDaysRu = []string{
		"воскресенье", "понедельник", "вторник",
		"среда", "четверг", "пятница", "суббота",
	}
)

// FmtDuration formats given duration like "HH:mm:ss" or "mm:ss".
func FmtDuration(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour
	mins := d / time.Minute
	secs := d % time.Minute / time.Second

	if hours == 0 {
		return fmt.Sprintf("%02d:%02d", mins, secs)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
}

// FmtRuDateFull formats given datetime like "04 июля 2026".
func FmtRuDateFull(t time.Time) string {
	return fmt.Sprintf("%02d %s %d", t.Day(), _monthsRu[t.Month()], t.Year())
}

// FmtRuMonth formats given datetime month like "июля".
func FmtRuMonth(t time.Time) string {
	return _monthsRu[t.Month()]
}

// FmtRuWeekDay formats given datetime week day like "суббота".
func FmtRuWeekDay(t time.Time) string {
	return _weekDaysRu[t.Weekday()]
}
