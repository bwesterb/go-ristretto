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

// Test vectors from https://ristretto.group/test_vectors/ristretto255.html
func TestRistretto255TestVectors(t *testing.T) {
	smallMultiples := []string{
		"0000000000000000000000000000000000000000000000000000000000000000",
		"e2f2ae0a6abc4e71a884a961c500515f58e30b6aa582dd8db6a65945e08d2d76",
		"6a493210f7499cd17fecb510ae0cea23a110e8d5b901f8acadd3095c73a3b919",
		"94741f5d5d52755ece4f23f044ee27d5d1ea1e2bd196b462166b16152a9d0259",
		"da80862773358b466ffadfe0b3293ab3d9fd53c5ea6c955358f568322daf6a57",
		"e882b131016b52c1d3337080187cf768423efccbb517bb495ab812c4160ff44e",
		"f64746d3c92b13050ed8d80236a7f0007c3b3f962f5ba793d19a601ebb1df403",
		"44f53520926ec81fbd5a387845beb7df85a96a24ece18738bdcfa6a7822a176d",
		"903293d8f2287ebe10e2374dc1a53e0bc887e592699f02d077d5263cdd55601c",
		"02622ace8f7303a31cafc63f8fc48fdc16e1c8c8d234b2f0d6685282a9076031",
		"20706fd788b2720a1ed2a5dad4952b01f413bcf0e7564de8cdc816689e2db95f",
		"bce83f8ba5dd2fa572864c24ba1810f9522bc6004afe95877ac73241cafdab42",
		"e4549ee16b9aa03099ca208c67adafcafa4c3f3e4e5303de6026e3ca8ff84460",
		"aa52e000df2e16f55fb1032fc33bc42742dad6bd5a8fc0be0167436c5948501f",
		"46376b80f409b29dc2b5f6f0c52591990896e5716f41477cd30085ab7f10301e",
		"e0c418f7c8d9c4cdd7395b93ea124f3ad99021bb681dfc3302a9d99a2e53e64e",
	}
	var B, pt ristretto.Point
	B.SetBase()
	pt.SetZero()
	for i, pt2 := range smallMultiples {
		if hex.EncodeToString(pt.Bytes()) != pt2 {
			t.Fatalf("%d * B = %s != %v", i, pt2, hex.EncodeToString(pt.Bytes()))
		}
		pt.Add(&B, &pt)
	}

	badEncodings := []string{
		"00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		// We ignore the top bit in a field element, but curve25519-dalek
		// does not. So in curve25519, these fail
		// "f3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		// "edffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
		"0100000000000000000000000000000000000000000000000000000000000000",
		"01ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f",
	}

	for _, ptHex := range badEncodings {
		var pt ristretto.Point
		var buf [32]byte
		tmp, _ := hex.DecodeString(ptHex)
		copy(buf[:], tmp)
		if pt.SetBytes(&buf) {
			t.Fatalf("%v should not decode", ptHex)
		}
	}

	// TODO add "encoded_hash_to_points" test
}
