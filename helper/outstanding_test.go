package helper

import (
	"progas-wms-be/enum"
	"testing"
	"time"
)

func TestDaysAtCustomer(t *testing.T) {
	past := func(hours int) *time.Time {
		ts := time.Now().Add(-time.Duration(hours) * time.Hour)
		return &ts
	}

	cases := []struct {
		name string
		in   *time.Time
		want int
	}{
		{"nil since", nil, 0},
		{"just arrived", past(1), 0},
		{"two days", past(50), 2},
		{"ten days", past(10 * 24), 10},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := DaysAtCustomer(c.in); got != c.want {
				t.Errorf("DaysAtCustomer() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestIsOverdueAtCustomer(t *testing.T) {
	past := func(hours int) *time.Time {
		ts := time.Now().Add(-time.Duration(hours) * time.Hour)
		return &ts
	}

	cases := []struct {
		name    string
		since   *time.Time
		maxDays int
		want    bool
	}{
		{"nil since is never overdue", nil, 30, false},
		{"zero max means no limit", past(400), 0, false},
		{"negative max means no limit", past(400), -1, false},
		{"under max", past(50), 5, false},
		{"exactly at max", past(50), 2, true},
		{"over max", past(10 * 24), 7, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsOverdueAtCustomer(c.since, c.maxDays); got != c.want {
				t.Errorf("IsOverdueAtCustomer(max=%v) = %v, want %v", c.maxDays, got, c.want)
			}
		})
	}
}

func TestCanSetMaxDaysAtCustomer(t *testing.T) {
	cases := []struct {
		status enum.CylinderStatus
		want   bool
	}{
		{enum.CylinderStatusEmpty, true},
		{enum.CylinderStatusReadyToFill, true},
		{enum.CylinderStatusFilled, true},
		{enum.CylinderStatusReady, true},
		{enum.CylinderStatusInTransit, true},
		{enum.CylinderStatusOutstanding, true},
		{enum.CylinderStatusMaintenance, true},
		{enum.CylinderStatusLost, false},
		{enum.CylinderStatusWriteOff, false},
	}
	for _, c := range cases {
		t.Run(string(c.status), func(t *testing.T) {
			if got := CanSetMaxDaysAtCustomer(c.status); got != c.want {
				t.Errorf("CanSetMaxDaysAtCustomer(%v) = %v, want %v", c.status, got, c.want)
			}
		})
	}
}
