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
