package service

import (
    "testing"

    "github.com/mhsanaei/3x-ui/v3/internal/database/model"
    "github.com/mhsanaei/3x-ui/v3/internal/xray"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestResolveWeightedClientTraffic(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    if err := db.AutoMigrate(
        &model.Inbound{},
        &model.ClientRecord{},
        &model.ClientInbound{},
    ); err != nil {
        t.Fatal(err)
    }

    inbound1 := model.Inbound{
        Port: 2085,
        Tag: "test-vless-2085",
        TrafficMultiplier: 2,
    }
    inbound2 := model.Inbound{
        Port: 443,
        Tag: "test-vless-443",
        TrafficMultiplier: 1,
    }

    if err := db.Create(&inbound1).Error; err != nil {
        t.Fatal(err)
    }
    if err := db.Create(&inbound2).Error; err != nil {
        t.Fatal(err)
    }

    client := model.ClientRecord{Email: "ali"}
    if err := db.Create(&client).Error; err != nil {
        t.Fatal(err)
    }

    for _, inbound := range []model.Inbound{inbound1, inbound2} {
        link := model.ClientInbound{
            ClientId: client.Id,
            InboundId: inbound.Id,
        }
        if err := db.Create(&link).Error; err != nil {
            t.Fatal(err)
        }
    }

    input := []*xray.ClientTraffic{
        {
            Email: TrafficStatsIdentity(inbound1.Id, "ali"),
            Up: 1000,
            Down: 2000,
        },
        {
            Email: TrafficStatsIdentity(inbound2.Id, "ali"),
            Up: 1000,
            Down: 2000,
        },
    }

    result, err := ResolveWeightedClientTraffic(db, input)
    if err != nil {
        t.Fatal(err)
    }

    var totalUp, totalDown int64
    for _, traffic := range result {
        if traffic.Email != "ali" {
            t.Fatalf("unexpected email: %s", traffic.Email)
        }
        totalUp += traffic.Up
        totalDown += traffic.Down
    }

    if totalUp != 3000 || totalDown != 6000 {
        t.Fatalf("weighted traffic = %d/%d, want 3000/6000",
            totalUp, totalDown)
    }
}
