package util

import (
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

// ParseDate parses a date string in YYYY-MM-DD format.
func ParseDate(s string) (time.Time, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: expected YYYY-MM-DD", s)
	}
	return t, nil
}

// FormatDate formats a time.Time to the YYYY-MM-DD string.
func FormatDate(t time.Time) string {
	return t.Format(dateLayout)
}

// IsValidDateRange returns true when date is within [minYearsBack, maxYearsAhead]
// relative to now.
func IsValidDateRange(date time.Time, minYearsBack, maxYearsAhead int) bool {
	now := time.Now()
	earliest := now.AddDate(-minYearsBack, 0, 0)
	latest := now.AddDate(maxYearsAhead, 0, 0)
	return !date.Before(earliest) && !date.After(latest)
}

// IsDateInFuture returns true when date is strictly after now (UTC day boundary).
func IsDateInFuture(date time.Time) bool {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	return date.UTC().Truncate(24 * time.Hour).After(today)
}

// IsDateInPast returns true when date is strictly before today (UTC day boundary).
func IsDateInPast(date time.Time) bool {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	return date.UTC().Truncate(24 * time.Hour).Before(today)
}

// MonthName returns the Indonesian month name for month 1..12.
func MonthName(month int) string {
	names := [...]string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month < 1 || month > 12 {
		return "Tidak diketahui"
	}
	return names[month-1]
}

// StartOfMonth returns the first day of the given month/year as a time.Time (UTC).
func StartOfMonth(year, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
}

// EndOfMonth returns the last moment of the given month/year.
func EndOfMonth(year, month int) time.Time {
	return StartOfMonth(year, month).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// DaysBetween returns the number of whole days between two dates (absolute value).
func DaysBetween(a, b time.Time) int {
	diff := b.Sub(a)
	if diff < 0 {
		diff = -diff
	}
	return int(diff.Hours() / 24)
}

// MonthsBetween returns the number of whole months from 'from' to 'to'.
// Returns 0 if 'to' is before 'from'.
func MonthsBetween(from, to time.Time) int {
	if to.Before(from) {
		return 0
	}
	years := to.Year() - from.Year()
	months := int(to.Month()) - int(from.Month())
	total := years*12 + months
	if total < 0 {
		return 0
	}
	return total
}
