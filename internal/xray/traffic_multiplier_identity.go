package xray

import (
    "encoding/base64"
    "strconv"
)

// TrafficStatsIdentity creates a distinct Xray stats email
// for one client on one inbound.
func TrafficStatsIdentity(inboundID int, email string) string {
    return "xuitm1_" +
        strconv.Itoa(inboundID) + "_" +
        base64.RawURLEncoding.EncodeToString([]byte(email))
}

// UsePerInboundStatsEmail currently enables separate
// statistics for VLESS inbounds with a non-default multiplier.
func UsePerInboundStatsEmail(protocol string, multiplier float64) bool {
    return protocol == "vless" &&
        multiplier > 0 &&
        multiplier != 1
}
