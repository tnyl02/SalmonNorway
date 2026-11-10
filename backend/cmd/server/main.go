// โปรแกรมหลังบ้านของระบบยืมอุปกรณ์ภาควิชา (CON-TECH-01 PostgreSQL, CON-TECH-02 รันใน Docker)
package main

import (
	"context"
	"log"
	"net/http"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/equipment"
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
	mux.Handle(equipment.PathPrefix, auth.Middleware(provider, eq.Routes()))

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatal(err)
	}
}
