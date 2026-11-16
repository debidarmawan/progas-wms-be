package mapper

import (
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"testing"
	"time"
)

func TestToCylinderResponseAging(t *testing.T) {
	daysAgo := func(days int) *time.Time {
		ts := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
		return &ts
	}

	cases := []struct {
		name        string
		status      enum.CylinderStatus
		since       *time.Time
		maxDays     int
		wantDays    int
		wantSince   bool
		wantOverdue bool
	}{
		{"still in warehouse", enum.CylinderStatusReady, nil, 30, 0, false, false},
		{"in transit day 2", enum.CylinderStatusInTransit, daysAgo(2), 30, 2, true, false},
		{"outstanding past limit", enum.CylinderStatusOutstanding, daysAgo(10), 7, 10, true, true},
		{"no limit never overdue", enum.CylinderStatusOutstanding, daysAgo(400), 0, 400, true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ToCylinderResponse(&model.Cylinder{Status: c.status, OutstandingSince: c.since, MaxDaysAtCustomer: c.maxDays})
			if got.DaysAtCustomer != c.wantDays {
				t.Errorf("DaysAtCustomer = %v, want %v", got.DaysAtCustomer, c.wantDays)
			}
			if (got.OutstandingSince != "") != c.wantSince {
				t.Errorf("OutstandingSince = %q, want set = %v", got.OutstandingSince, c.wantSince)
			}
			if got.IsOverdue != c.wantOverdue {
				t.Errorf("IsOverdue = %v, want %v", got.IsOverdue, c.wantOverdue)
			}
			if got.MaxDaysAtCustomer != c.maxDays {
				t.Errorf("MaxDaysAtCustomer = %v, want %v", got.MaxDaysAtCustomer, c.maxDays)
			}
		})
	}
}
