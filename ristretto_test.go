package ristretto_test

import (
	"encoding/hex"
	"testing"

	"github.com/bwesterb/go-ristretto"
)

func TestPointDerive(t *testing.T) {
	testVectors := []struct{ in, out string }{
		{"test", "b01d60504aa5f4c5bd9a7541c457661f9a789d18cb4e136e91d3c953488bd208"},
		{"pep", "3286c8d171dec02e70549c280d62524430408a781efc07e4428d1735671d195b"},
		{"ristretto", "c2f6bb4c4dab8feab66eab09e77e79b36095c86b3cd1145b9a2703205858d712"},
		{"elligator", "784c727b1e8099eb94e5a8edbd260363567fdbd35106a7a29c8b809cd108b322"},
	}
	for _, v := range testVectors {
		var p ristretto.Point
		p.Derive([]byte(v.in))
		out2 := hex.EncodeToString(p.Bytes())
		if out2 != v.out {
			t.Fatalf("Derive(%v) = %v != %v", v.in, v.out, out2)
		}
	}
}
