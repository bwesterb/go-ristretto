package ristretto

// A number modulo the prime l, where l is the order of the Ristretto group
// over Edwards25519.
//
// The scalar s is represented as an array s[0], ... s[31] with 0 <= s[i] <= 255
// and s = s[0] + s[1] * 256 + s[2] * 65536 + ... + s[31] * 256^31.
// We use uint32 (instead of uint8) so that we have some spare room during
// computations.
type Scalar [32]uint32
