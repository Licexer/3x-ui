package service

import "testing"

func TestWeightedTraffic(t *testing.T) {
    cases := []struct {
        name string
        bytes int64
        multiplier float64
        want int64
        wantErr bool
    }{
        {"default", 1000, 0, 1000, false},
        {"normal", 1000, 1, 1000, false},
        {"double", 1000, 2, 2000, false},
        {"half", 1000, 0.5, 500, false},
        {"negative bytes", -1, 2, 0, true},
        {"negative factor", 1000, -2, 0, true},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := WeightedTraffic(tc.bytes, tc.multiplier)
            if (err != nil) != tc.wantErr {
                t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
            }
            if !tc.wantErr && got != tc.want {
                t.Fatalf("got %d, want %d", got, tc.want)
            }
        })
    }
}
