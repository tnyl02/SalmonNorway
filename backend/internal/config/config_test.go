package config

import (
	"slices"
	"testing"
	"time"
)

// 006 T-02 เสร็จเมื่อ: ค่าของการรับคืนตรงกับ spec (FR-RET-05, ASM-06)
func TestReturnPhotoLimits(t *testing.T) {
	if ReturnPhotoMin != 1 || ReturnPhotoMax != 5 {
		t.Fatalf("จำนวนรูปต้องเป็น 1 ถึง 5 ได้ %d ถึง %d", ReturnPhotoMin, ReturnPhotoMax)
	}
	if ReturnPhotoMaxBytes != 5*1024*1024 {
		t.Fatalf("ขนาดต่อรูปต้องเป็น 5 MB ได้ %d ไบต์", ReturnPhotoMaxBytes)
	}
	if ReturnRequestMaxBytes <= ReturnPhotoMax*ReturnPhotoMaxBytes {
		t.Fatal("เพดาน request ต้องใหญ่กว่ารูปครบจำนวนสูงสุด ไม่งั้นรูป 5 รูปขนาดเต็มจะถูกตัด")
	}
	if !slices.Equal(ReturnPhotoTypes, []string{"image/jpeg", "image/png"}) {
		t.Fatalf("ชนิดรูปต้องเป็น JPG/PNG ได้ %v", ReturnPhotoTypes)
	}
}

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
