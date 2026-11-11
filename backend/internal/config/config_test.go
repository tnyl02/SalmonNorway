package config

import (
	"testing"
	"time"
)

// T-03 เสร็จเมื่อ: สถานะที่ยกเลิกได้มี 3 ค่า (CON-CXL-02) และ CancelDBTimeout น้อยกว่า 15 วินาที (NFR-PERF-CXL-01, ASM-06)
func TestConfig_CancellableStatuses(t *testing.T) {
	got := CancellableStatuses()
	want := []string{"requested", "approved_by_officer", "pending_advisor_approval"}
	if len(got) != len(want) {
		t.Fatalf("CancellableStatuses() = %v ต้องมี %d ค่า: %v", got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ค่าที่ %d = %q ต้องเป็น %q", i, got[i], want[i])
		}
	}

	// แก้ slice ที่ได้ไปแล้ว ค่ากลางต้องไม่เปลี่ยน
	got[0] = "approved"
	if CancellableStatuses()[0] != "requested" {
		t.Errorf("CancellableStatuses() ต้องคืน slice ใหม่ทุกครั้ง")
	}
}

func TestConfig_CancelLimits(t *testing.T) {
	if CancelReasonMaxLen != 200 {
		t.Errorf("CancelReasonMaxLen = %d ต้องเป็น 200 (FR-CXL-08)", CancelReasonMaxLen)
	}
	if CancelPerfP95 != 15*time.Second {
		t.Errorf("CancelPerfP95 = %v ต้องเป็น 15s (NFR-PERF-CXL-01)", CancelPerfP95)
	}
	if CancelDBTimeout >= CancelPerfP95 {
		t.Errorf("CancelDBTimeout = %v ต้องน้อยกว่า CancelPerfP95 = %v", CancelDBTimeout, CancelPerfP95)
	}
}
