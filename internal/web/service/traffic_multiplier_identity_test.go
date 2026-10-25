package service

import "testing"

func TestTrafficStatsIdentityRoundTrip(t *testing.T) {
    email := "ali@example.com"

    identity := TrafficStatsIdentity(12, email)

    id, decoded, ok := ParseTrafficStatsIdentity(identity)
    if !ok {
        t.Fatal("identity not decoded")
    }
    if id != 12 {
        t.Fatalf("inbound ID = %d, want 12", id)
    }
    if decoded != email {
        t.Fatalf("email = %q, want %q", decoded, email)
    }

    other := TrafficStatsIdentity(13, email)
    if identity == other {
        t.Fatal("different inbounds generated identical identities")
    }

    _, _, ok = ParseTrafficStatsIdentity(email)
    if ok {
        t.Fatal("legacy email must not parse as stats identity")
    }
}

func TestTrafficStatsIdentityInvalidInputs(t *testing.T) {
    invalid := []string{
        "",
        "alice",
        "xuitm1_0_YWxpY2U",
        "xuitm1_-1_YWxpY2U",
        "xuitm1_12_",
        "xuitm1_12_!",
        "xuitm1_012_YWxpY2U",
        "xuitm1_12_YS9i",
        "xuitm1_12_YWxpY2U=",
    }

    for _, identity := range invalid {
        if _, _, ok := ParseTrafficStatsIdentity(identity); ok {
            t.Errorf("invalid identity %q was accepted", identity)
        }
    }
}
