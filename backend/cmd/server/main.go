// โปรแกรมหลังบ้านของระบบยืมอุปกรณ์ภาควิชา (CON-TECH-01 PostgreSQL, CON-TECH-02 รันใน Docker)
package main

import (
	"context"
	"log"

	"salmonnorway/backend/internal/approve"
	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/borrowrequest"
	"salmonnorway/backend/internal/cancelrequest"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/equipment"
	"salmonnorway/backend/internal/handover"
	"salmonnorway/backend/internal/platform"
	"salmonnorway/backend/internal/returns"
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
	// path ลูกของฟีเจอร์นี้ เช่น /v1/loan-requests/{id}/attachment (T-15, FR-BRQ-10 ข้อ ก)
	// mux ชั้นนอกจับ path แบบเป๊ะ ๆ ถ้าไม่ผูกตัวที่ลงท้ายด้วย / ไว้ด้วย path ลูกจะไปไม่ถึง router ของฟีเจอร์
	mux.Handle(borrowrequest.PathPrefix+"/", auth.Middleware(provider, br.Routes()))
	// ลำดับคำขอที่ค้างอยู่ของอุปกรณ์ชิ้นนั้น เจ้าของคือ 002 แม้ path จะอยู่ใต้ /v1/equipment/ ของ 001
	// (สัญญา API ใน specs/changes/CR-005.md, T-18 ของ 002) Go 1.22 ส่งให้ pattern ที่เจาะจงกว่า จึงไม่ชนกัน
	mux.Handle("GET "+borrowrequest.QueuePath, auth.Middleware(provider, br.Routes()))

	// ฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (T-05) ต้องยืนยันตัวตนและเป็นเจ้าหน้าที่ (IF-IDP-01, NFR-SEC-01, AC-APV-05)
	ap := approve.NewHandler(approve.NewService(approve.NewRepository(db)))
	mux.Handle(approve.PathPrefix, auth.Middleware(provider, ap.Routes()))
	// 004 งานเบื้องหลัง (T-30) แก้ไฟล์กลางเพราะ 004:
	//   - ตัวส่ง Email แจ้งผล ภายใน 60 วินาที (FR-APV-02) ใช้ตัวส่งและที่อยู่จำลองไปก่อน (IF-NOT-01, รอ Q-APV-04)
	//   - ปฏิเสธคำขอยืมปกติที่ค้างตอนเที่ยงคืน และไล่ปฏิเสธคำขอที่เลยเวลาทันทีตอนเริ่มทำงาน (FR-APV-01 ฉ, Q-APV-21)
	go approve.NewSender(db, approve.LogMailer{}, approve.MockRecipientResolver{}).Run(ctx)
	go approve.NewAutoRejecter(db).Run(ctx)

	// 003 ยกเลิกคำขอยืม (T-04) หลัง auth.Middleware เหมือน 001 (IF-IDP-01) แก้ไฟล์กลางเพราะ 003
	cx := cancelrequest.NewHandler(cancelrequest.NewService(cancelrequest.NewRepository(db)))
	cxRoutes := auth.Middleware(provider, cx.Routes())
	mux.Handle(cancelrequest.PathPrefixMine, cxRoutes)
	mux.Handle(cancelrequest.PathPrefixRequests, cxRoutes)

	// 006 บันทึกการรับคืน (specs/006-return/plan.md ข้อ 4) เฉพาะเจ้าหน้าที่ภาควิชาเดียวกัน (FR-RET-01, AC-RET-08)
	// Routes ห่อ auth.Middleware และตรวจบทบาทเองทุก route (IF-IDP-01) แก้ไฟล์กลางเพราะ 006
	returns.NewHandler(returns.NewRepository(db)).Routes(provider)(mux)

	// ทุก response ห้าม cache เพื่อให้รีเฟรชแล้วเห็นสถานะล่าสุด (NFR-DATA-01, T-16)
	// 405 ของทุกฟีเจอร์ตอบ {code: METHOD_NOT_ALLOWED, message} (CON-TECH-03, T-40 ของ 001 แก้ไฟล์กลาง แจ้งทีมแล้ว 2569-11-12)
	// HTTPS เท่านั้นเมื่อตั้ง TLS_CERT_FILE และ TLS_KEY_FILE ด้วย TLS 1.2 ขึ้นไป (CON-TECH-02, NFR-SEC-01, AC-APV-19)
	// 004 รับทำเอง Q-APV-38 (T-50 ของ 004 แก้ไฟล์กลาง) ไม่ตั้ง = HTTP แบบเดิมสำหรับเครื่องพัฒนา
	handler := platform.NoStore(platform.MethodNotAllowedJSON(mux))
	if err := platform.ListenAndServe(cfg.HTTPAddr, cfg.TLSCertFile, cfg.TLSKeyFile, handler); err != nil {
		log.Fatal(err)
	}
}
