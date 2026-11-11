// โปรแกรมหลังบ้านของระบบยืมอุปกรณ์ภาควิชา (CON-TECH-01 PostgreSQL, CON-TECH-02 รันใน Docker)
package main

import (
	"context"
	"log"
	"net/http"

	"salmonnorway/backend/internal/approve"
	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/borrowrequest"
	"salmonnorway/backend/internal/cancelrequest"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/equipment"
	"salmonnorway/backend/internal/handover"
	"salmonnorway/backend/internal/platform"
)

// path ของไฟล์ migration และข้อมูลตั้งต้น เทียบจากโฟลเดอร์ backend/ (รันด้วย `cd backend && go run ./cmd/server`)
const (
	migrationsDir = "db/migrations"
	seedFile      = "db/seed/seed.sql"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// สร้างหรืออัปเดตตารางก่อนเปิดรับ request (T-02, T-31)
	if err := platform.Migrate(cfg.DatabaseURL, migrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	db, err := platform.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database config: %v", err)
	}
	defer db.Close()

	// ใส่ข้อมูลตั้งต้นเฉพาะตอนฐานข้อมูลยังว่าง (ASM-01 ข้อมูลตัวอย่างสำหรับเครื่องพัฒนา)
	// ข้อมูลตั้งต้นใช้ id ตายตัว ใส่ซ้ำไม่ได้ จึงเช็กก่อนทุกครั้งที่เปิดหลังบ้าน
	var hasData bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM department)`).Scan(&hasData); err != nil {
		log.Fatalf("check seed: %v", err)
	}
	if !hasData {
		if err := platform.Seed(ctx, cfg.DatabaseURL, seedFile); err != nil {
			log.Fatalf("seed: %v", err)
		}
		log.Printf("seeded demo data from %s", seedFile)
	}

	// ทุก endpoint ของฟีเจอร์ต้องผ่านการยืนยันตัวตนก่อน (IF-IDP-01, NFR-SEC-01)
	// ใช้ระบบยืนยันตัวตนจำลองไปก่อน รอ Q-SH-03 (T-30)
	provider := auth.NewMockProvider()

	mux := platform.NewRouter(db) // /healthz ไม่ต้องยืนยันตัวตน ใช้ตรวจ container
	eq := equipment.NewHandler(equipment.NewService(equipment.NewRepository(db)))
	ho := handover.NewHandler(db, provider) // 005 บันทึกการส่งมอบ เฉพาะเจ้าหน้าที่ (CON-HND-01)
	mux.Handle(equipment.PathPrefix, auth.Middleware(provider, eq.Routes()))
	mux.Handle(handover.PathPrefix, auth.Middleware(provider, ho.Routes()))

	// ฟีเจอร์ 002 ส่งคำขอยืม (T-05) ต้องผ่านการยืนยันตัวตนเช่นกัน (IF-IDP-01, NFR-SEC-01)
	br := borrowrequest.NewHandler(borrowrequest.NewService(borrowrequest.NewRepository(db)))
	mux.Handle(borrowrequest.PathPrefix, auth.Middleware(provider, br.Routes()))

	// ฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (T-05) ต้องยืนยันตัวตนและเป็นเจ้าหน้าที่ (IF-IDP-01, NFR-SEC-01, AC-APV-05)
	ap := approve.NewHandler(approve.NewService(approve.NewRepository(db)))
	mux.Handle(approve.PathPrefix, auth.Middleware(provider, ap.Routes()))

	// 003 ยกเลิกคำขอยืม (T-04) หลัง auth.Middleware เหมือน 001 (IF-IDP-01) แก้ไฟล์กลางเพราะ 003
	cx := cancelrequest.NewHandler(cancelrequest.NewService(cancelrequest.NewRepository(db)))
	cxRoutes := auth.Middleware(provider, cx.Routes())
	mux.Handle(cancelrequest.PathPrefixMine, cxRoutes)
	mux.Handle(cancelrequest.PathPrefixRequests, cxRoutes)

	log.Printf("listening on %s", cfg.HTTPAddr)
	// ทุก response ห้าม cache เพื่อให้รีเฟรชแล้วเห็นสถานะล่าสุด (NFR-DATA-01, T-16)
	if err := http.ListenAndServe(cfg.HTTPAddr, platform.NoStore(mux)); err != nil {
		log.Fatal(err)
	}
}
