package service

import (
    "encoding/base64"
    "fmt"
    "strconv"
    "strings"
)

const trafficIdentityPrefix = "xuitm1_"

// TrafficStatsIdentity returns a unique stats identity per
// (inbound ID, original email) pair.
func TrafficStatsIdentity(inboundID int, email string) string {
    encoded := base64.RawURLEncoding.EncodeToString([]byte(email))
    return fmt.Sprintf("%s%d_%s", trafficIdentityPrefix, inboundID, encoded)
}

// ParseTrafficStatsIdentity decodes a generated stats identity.
// The second return value is false for legacy identities.
func ParseTrafficStatsIdentity(identity string) (int, string, bool) {
    if !strings.HasPrefix(identity, trafficIdentityPrefix) {
        return 0, "", false
    }

    tail := strings.TrimPrefix(identity, trafficIdentityPrefix)
    parts := strings.SplitN(tail, "_", 2)
    if len(parts) != 2 {
        return 0, "", false
    }

    inboundID, err := strconv.Atoi(parts[0])
    if err != nil || inboundID <= 0 {
        return 0, "", false
    }

    raw, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil || len(raw) == 0 {
        return 0, "", false
    }

    email := string(raw)
    if TrafficStatsIdentity(inboundID, email) != identity {
        return 0, "", false
    }
    if validateClientEmail(email) != nil {
        return 0, "", false
    }

    return inboundID, email, true
}
