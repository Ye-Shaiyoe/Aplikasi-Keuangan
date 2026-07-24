package util_test

import (
	"testing"
	"time"

	"github.com/akrom/finance-backend/internal/util"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2024-01-15", false},
		{"2025-12-31", false},
		{"invalid", true},
		{"2024/01/15", true},
		{"15-01-2024", true},
	}

	for _, tt := range tests {
		got, err := util.ParseDate(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseDate(%q) err = %v, wantErr = %v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && util.FormatDate(got) != tt.input {
			t.Errorf("FormatDate(ParseDate(%q)) = %q", tt.input, util.FormatDate(got))
		}
	}
}

func TestMonthName(t *testing.T) {
	if util.MonthName(1) != "Januari" {
		t.Errorf("MonthName(1) = %q, want 'Januari'", util.MonthName(1))
	}
	if util.MonthName(12) != "Desember" {
		t.Errorf("MonthName(12) = %q, want 'Desember'", util.MonthName(12))
	}
	if util.MonthName(0) != "Tidak diketahui" {
		t.Errorf("MonthName(0) = %q, want 'Tidak diketahui'", util.MonthName(0))
	}
}

func TestStartAndEndOfMonth(t *testing.T) {
	start := util.StartOfMonth(2024, 2) // Feb 2024 (leap year)
	if start.Day() != 1 || start.Month() != time.February {
		t.Errorf("StartOfMonth = %v", start)
	}

	end := util.EndOfMonth(2024, 2)
	if end.Day() != 29 {
		t.Errorf("EndOfMonth Feb 2024 day = %d, want 29", end.Day())
	}
}

func TestDaysBetween(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	if util.DaysBetween(t1, t2) != 9 {
		t.Errorf("DaysBetween = %d, want 9", util.DaysBetween(t1, t2))
	}
	if util.DaysBetween(t2, t1) != 9 {
		t.Errorf("DaysBetween reverse = %d, want 9", util.DaysBetween(t2, t1))
	}
}

func TestMonthsBetween(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	if util.MonthsBetween(t1, t2) != 5 {
		t.Errorf("MonthsBetween = %d, want 5", util.MonthsBetween(t1, t2))
	}
	if util.MonthsBetween(t2, t1) != 0 {
		t.Errorf("MonthsBetween past = %d, want 0", util.MonthsBetween(t2, t1))
	}
}
