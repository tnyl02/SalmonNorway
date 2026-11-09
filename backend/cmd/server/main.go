// โปรแกรมหลังบ้านของระบบยืมอุปกรณ์ภาควิชา (CON-TECH-01 PostgreSQL, CON-TECH-02 รันใน Docker)
package main

import (
	"context"
	"log"
	"net/http"

	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

func main() {
	cfg := config.Load()

	db, err := platform.OpenDB(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database config: %v", err)
	}
	defer db.Close()

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, platform.NewRouter(db)); err != nil {
		log.Fatal(err)
	}
}
