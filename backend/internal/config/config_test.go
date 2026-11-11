package config

import (
	"testing"
	"time"
)

// 006 T-02 เสร็จเมื่อ: สภาพมี 2 ค่า และสถานะหลังคืนเป็นค่าใน CON-DATA-01 (FR-RET-02, FR-RET-03, FR-RET-04)
func TestReturnConditionStatus(t *testing.T) {
	want := map[string]string{"normal": "available", "damaged": "unavailable"}
	if len(ReturnConditionStatus) != len(want) {
		t.Fatalf("สภาพต้องมี 2 ค่า ได้ %v", ReturnConditionStatus)
	}
	for c, s := range want {
		if ReturnConditionStatus[c] != s {
			t.Errorf("สภาพ %q ต้องได้สถานะ %q ได้ %q", c, s, ReturnConditionStatus[c])
		}
	}
	if _, ok := ReturnConditionStatus["lost"]; ok {
		t.Error("สูญหายไม่ใช่สภาพ (FR-RET-02)")
	}
}

// เขตเวลาต้องโหลดได้จริง (glossary.md)
func TestTimezoneLoads(t *testing.T) {
	if _, err := time.LoadLocation(Timezone); err != nil {
		t.Fatalf("โหลดเขตเวลา %q ไม่ได้: %v", Timezone, err)
	}
}
