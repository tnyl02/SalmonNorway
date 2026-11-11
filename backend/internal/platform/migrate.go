package platform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // driver pgx5:// ของ golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

// Migrate สร้างหรืออัปเดตตารางจากไฟล์ใน migrationsDir ด้วย golang-migrate (plan.md ข้อ 2, CON-TECH-01)
// เรียกซ้ำได้ ถ้าเป็นเวอร์ชันล่าสุดแล้วจะไม่ทำอะไร
func Migrate(databaseURL, migrationsDir string) error {
	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}
	// ToSlash: บน Windows abs เป็น D:\... ซึ่ง URL อ่านไม่ได้ (invalid port) เปลี่ยนเป็น D:/... บน Linux ไม่เปลี่ยนอะไร
	// แก้โดย 003 (2569-11-11) ไฟล์ของกลางของ 001 ต้องแจ้งเจ้าของ 001
	m, err := migrate.New("file://"+filepath.ToSlash(abs), toPgx5URL(databaseURL))
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// Seed ใส่ข้อมูลตั้งต้นจากไฟล์ SQL สำหรับเครื่องพัฒนา (ASM-01)
func Seed(ctx context.Context, databaseURL, seedFile string) error {
	sql, err := os.ReadFile(seedFile)
	if err != nil {
		return err
	}
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, string(sql))
	return err
}

// golang-migrate เลือก driver จาก scheme ของ URL จึงเปลี่ยน postgres:// เป็น pgx5://
func toPgx5URL(databaseURL string) string {
	for _, p := range []string{"postgresql://", "postgres://"} {
		if strings.HasPrefix(databaseURL, p) {
			return "pgx5://" + strings.TrimPrefix(databaseURL, p)
		}
	}
	return databaseURL
}
