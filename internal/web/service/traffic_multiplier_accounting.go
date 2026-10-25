package service

import (
    "fmt"

    "github.com/mhsanaei/3x-ui/v3/internal/database/model"
    "github.com/mhsanaei/3x-ui/v3/internal/xray"
    "gorm.io/gorm"
)

// ResolveWeightedClientTraffic converts per-inbound Xray statistics
// into billable traffic under the original shared client email.
func ResolveWeightedClientTraffic(
    tx *gorm.DB,
    traffics []*xray.ClientTraffic,
) ([]*xray.ClientTraffic, error) {
    result := make([]*xray.ClientTraffic, 0, len(traffics))

    for _, traffic := range traffics {
        if traffic == nil {
            continue
        }

        inboundID, email, encoded := ParseTrafficStatsIdentity(traffic.Email)
        if !encoded {
            result = append(result, traffic)
            continue
        }

        var inbound model.Inbound
        if err := tx.Select("id", "traffic_multiplier").
            Where("id = ?", inboundID).
            Take(&inbound).Error; err != nil {
            return nil, fmt.Errorf("resolve traffic inbound %d: %w", inboundID, err)
        }

        var count int64
        if err := tx.Table("client_inbounds").
            Joins("JOIN clients ON clients.id = client_inbounds.client_id").
            Where("client_inbounds.inbound_id = ? AND clients.email = ?", inboundID, email).
            Count(&count).Error; err != nil {
            return nil, fmt.Errorf("verify traffic ownership: %w", err)
        }

        if count != 1 {
            return nil, fmt.Errorf("traffic identity has no unique owner: inbound=%d email=%q", inboundID, email)
        }

        up, err := WeightedTraffic(traffic.Up, inbound.TrafficMultiplier)
        if err != nil {
            return nil, err
        }

        down, err := WeightedTraffic(traffic.Down, inbound.TrafficMultiplier)
        if err != nil {
            return nil, err
        }

        copy := *traffic
        copy.Email = email
        copy.Up = up
        copy.Down = down

        result = append(result, &copy)
    }

    return result, nil
}
