// โปรแกรมหลังบ้านของระบบยืมอุปกรณ์ภาควิชา (CON-TECH-01 PostgreSQL, CON-TECH-02 รันใน Docker)
package main

import (
	"context"
	"log"
	"net/http"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
	"salmonnorway/backend/internal/returns"
)

func main() {
	cfg := config.Load()

	db, err := platform.OpenDB(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database config: %v", err)
	}
	defer db.Close()

	// ระบบยืนยันตัวตนจำลอง (IF-IDP-01 รอ Q-SH-03) ใช้ร่วมทุกฟีเจอร์
	provider := auth.NewMockProvider()

	// 006-return บันทึกการรับคืน (specs/006-return/plan.md ข้อ 4)
	returnsHandler := returns.NewHandler(returns.NewRepository(db))

	router := platform.NewRouter(db,
		platform.Routes(returnsHandler.Routes(provider)),
	)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
