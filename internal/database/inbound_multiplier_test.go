package database

import (
    "testing"

    "github.com/mhsanaei/3x-ui/v3/internal/database/model"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestInboundTrafficMultiplierMigration(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    if err := db.AutoMigrate(&model.Inbound{}); err != nil {
        t.Fatal(err)
    }

    if !db.Migrator().HasColumn(&model.Inbound{}, "traffic_multiplier") {
        t.Fatal("traffic_multiplier column missing")
    }

    if err := db.Exec(
        "INSERT INTO inbounds (port, tag, protocol) VALUES (?, ?, ?)",
        2085, "test-2085", "vless",
    ).Error; err != nil {
        t.Fatal(err)
    }

    var factor float64
    if err := db.Raw(
        "SELECT traffic_multiplier FROM inbounds WHERE port = ?",
        2085,
    ).Scan(&factor).Error; err != nil {
        t.Fatal(err)
    }

    if factor != 1 {
        t.Fatalf("default multiplier = %v, want 1", factor)
    }
}
