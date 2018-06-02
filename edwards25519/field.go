package edwards25519

import (
	// Requires for FieldElement.[Set]BigInt().  Obviously not used for actual
	// implementation, as operations on big.Ints are  not constant-time.
	"math/big"
)

// Element of the field GF(2^255 - 19) over which the elliptic
// curve Edwards25519 is defined.
type FieldElement [10]int32

var (
	feZero     FieldElement
	feOne      = FieldElement{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	feMinusOne = FieldElement{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// sqrt(-1)
	feI = FieldElement{
		-32595792, -7943725, 9377950, 3500415, 12389472,
		-272473, -25146209, -2005654, 326686, 11406482,
	}
)

// Sets fe to -a. Returns fe.
func (fe *FieldElement) Neg(a *FieldElement) *FieldElement {
	fe[0] = -a[0]
	fe[1] = -a[1]
	fe[2] = -a[2]
	fe[3] = -a[3]
	fe[4] = -a[4]
	fe[5] = -a[5]
	fe[6] = -a[6]
	fe[7] = -a[7]
	fe[8] = -a[8]
	fe[9] = -a[9]
	return fe
}

// Set fe to i, the root of -1.  Returns fe.
func (fe *FieldElement) SetI() *FieldElement {
	copy(fe[:], feI[:])
	return fe
}

// Set fe to 0.  Returns fe.
func (fe *FieldElement) SetZero() *FieldElement {
	copy(fe[:], feZero[:])
	return fe
}

// Set fe to 1.  Returns fe.
func (fe *FieldElement) SetOne() *FieldElement {
	copy(fe[:], feOne[:])
	return fe
}

// Sets fe to 2*a without normalizing.  Returns fe.
func (fe *FieldElement) double(a *FieldElement) *FieldElement {
	return fe.add(a, a)
}

// Sets fe to a + b without normalizing.  Returns fe.
func (fe *FieldElement) add(a, b *FieldElement) *FieldElement {
	fe[0] = a[0] + b[0]
	fe[1] = a[1] + b[1]
	fe[2] = a[2] + b[2]
	fe[3] = a[3] + b[3]
	fe[4] = a[4] + b[4]
	fe[5] = a[5] + b[5]
	fe[6] = a[6] + b[6]
	fe[7] = a[7] + b[7]
	fe[8] = a[8] + b[8]
	fe[9] = a[9] + b[9]
	return fe
}

// Sets fe to a + b.  Returns fe.
func (fe *FieldElement) Add(a, b *FieldElement) *FieldElement {
	return fe.add(a, b).normalize()
}

// Sets fe to a - b without normalizing.  Returns fe.
func (fe *FieldElement) sub(a, b *FieldElement) *FieldElement {
	fe[0] = a[0] - b[0]
	fe[1] = a[1] - b[1]
	fe[2] = a[2] - b[2]
	fe[3] = a[3] - b[3]
	fe[4] = a[4] - b[4]
	fe[5] = a[5] - b[5]
	fe[6] = a[6] - b[6]
	fe[7] = a[7] - b[7]
	fe[8] = a[8] - b[8]
	fe[9] = a[9] - b[9]
	return fe
}

// Sets fe to a - b.  Returns fe.
func (fe *FieldElement) Sub(a, b *FieldElement) *FieldElement {
	return fe.sub(a, b).normalize()
}

// Sets fe to a.  Returns fe.
func (fe *FieldElement) Set(a *FieldElement) *FieldElement {
	copy(fe[:], a[:])
	return fe
}

// Returns 1 if b < 0 and 0 otherwise.
func negative(b int32) int32 {
	return (b >> 31) & 1
}

// Returns 1 if b == c and 0 otherwise.  Assumes 2^15  <= b, c < 2^30.
func equal15(b, c int32) int32 {
	ub := uint16(b)
	uc := uint16(c)
	x := uint32(ub ^ uc)
	x--
	return int32(x >> 31)
}

// Returns 1 if b == c and 0 otherwise.  Assumes 0 <= b, c < 2^30.
func equal30(b, c int32) int32 {
	x := uint32(b ^ c)
	x--
	return int32(x >> 31)
}

// Interprets a 3-byte unsigned little endian byte-slice as int64
func load3(in []byte) int64 {
	var r int64
	r = int64(in[0])
	r |= int64(in[1]) << 8
	r |= int64(in[2]) << 16
	return r
}

// Interprets a 4-byte unsigned little endian byte-slice as int64
func load4(in []byte) int64 {
	var r int64
	r = int64(in[0])
	r |= int64(in[1]) << 8
	r |= int64(in[2]) << 16
	r |= int64(in[3]) << 24
	return r
}

// Reduce the even coefficients to below 1.01*2^25 and the odd coefficients
// to below 1.01*2^24.  Returns fe.
func (fe *FieldElement) normalize() *FieldElement {
	return fe.setReduced(
		int64(fe[0]), int64(fe[1]), int64(fe[2]), int64(fe[3]), int64(fe[4]),
		int64(fe[5]), int64(fe[6]), int64(fe[7]), int64(fe[8]), int64(fe[9]))
}

// Set fe to h0 + h1*2^26 + h2*2^51 + ... + h9*2^230.  Requires a little
// headroom in the inputs to store the carries.  Returns fe.
func (fe *FieldElement) setReduced(
	h0, h1, h2, h3, h4, h5, h6, h7, h8, h9 int64) *FieldElement {
	var c0, c1, c2, c3, c4, c5, c6, c7, c8, c9 int64

	c0 = (h0 + (1 << 25)) >> 26
	h1 += c0
	h0 -= c0 << 26
	c4 = (h4 + (1 << 25)) >> 26
	h5 += c4
	h4 -= c4 << 26

	c1 = (h1 + (1 << 24)) >> 25
	h2 += c1
	h1 -= c1 << 25
	c5 = (h5 + (1 << 24)) >> 25
	h6 += c5
	h5 -= c5 << 25

	c2 = (h2 + (1 << 25)) >> 26
	h3 += c2
	h2 -= c2 << 26
	c6 = (h6 + (1 << 25)) >> 26
	h7 += c6
	h6 -= c6 << 26

	c3 = (h3 + (1 << 24)) >> 25
	h4 += c3
	h3 -= c3 << 25
	c7 = (h7 + (1 << 24)) >> 25
	h8 += c7
	h7 -= c7 << 25

	c4 = (h4 + (1 << 25)) >> 26
	h5 += c4
	h4 -= c4 << 26
	c8 = (h8 + (1 << 25)) >> 26
	h9 += c8
	h8 -= c8 << 26

	c9 = (h9 + (1 << 24)) >> 25
	h0 += c9 * 19
	h9 -= c9 << 25

	c0 = (h0 + (1 << 25)) >> 26
	h1 += c0
	h0 -= c0 << 26

	fe[0] = int32(h0)
	fe[1] = int32(h1)
	fe[2] = int32(h2)
	fe[3] = int32(h3)
	fe[4] = int32(h4)
	fe[5] = int32(h5)
	fe[6] = int32(h6)
	fe[7] = int32(h7)
	fe[8] = int32(h8)
	fe[9] = int32(h9)

	return fe
}

// Set fe to a if b == 1.  Requires b to be either 0 or 1.
func (fe *FieldElement) ConditionalSet(a *FieldElement, b int32) {
	b = -b // b == 0b11111111111111111111111111111111 or 0.
	fe[0] ^= b & (fe[0] ^ a[0])
	fe[1] ^= b & (fe[1] ^ a[1])
	fe[2] ^= b & (fe[2] ^ a[2])
	fe[3] ^= b & (fe[3] ^ a[3])
	fe[4] ^= b & (fe[4] ^ a[4])
	fe[5] ^= b & (fe[5] ^ a[5])
	fe[6] ^= b & (fe[6] ^ a[6])
	fe[7] ^= b & (fe[7] ^ a[7])
	fe[8] ^= b & (fe[8] ^ a[8])
	fe[9] ^= b & (fe[9] ^ a[9])
}

// Returns little endian representation of fe.
func (fe *FieldElement) Bytes() [32]byte {
	var ret [32]byte
	fe.BytesInto(&ret)
	return ret
}

// Write fe to s in little endian.  Returns fe.
func (fe *FieldElement) BytesInto(s *[32]byte) *FieldElement {
	var carry [10]int32

	q := (19*fe[9] + (1 << 24)) >> 25
	q = (fe[0] + q) >> 26
	q = (fe[1] + q) >> 25
	q = (fe[2] + q) >> 26
	q = (fe[3] + q) >> 25
	q = (fe[4] + q) >> 26
	q = (fe[5] + q) >> 25
	q = (fe[6] + q) >> 26
	q = (fe[7] + q) >> 25
	q = (fe[8] + q) >> 26
	q = (fe[9] + q) >> 25

	fe[0] += 19 * q

	carry[0] = fe[0] >> 26
	fe[1] += carry[0]
	fe[0] -= carry[0] << 26
	carry[1] = fe[1] >> 25
	fe[2] += carry[1]
	fe[1] -= carry[1] << 25
	carry[2] = fe[2] >> 26
	fe[3] += carry[2]
	fe[2] -= carry[2] << 26
	carry[3] = fe[3] >> 25
	fe[4] += carry[3]
	fe[3] -= carry[3] << 25
	carry[4] = fe[4] >> 26
	fe[5] += carry[4]
	fe[4] -= carry[4] << 26
	carry[5] = fe[5] >> 25
	fe[6] += carry[5]
	fe[5] -= carry[5] << 25
	carry[6] = fe[6] >> 26
	fe[7] += carry[6]
	fe[6] -= carry[6] << 26
	carry[7] = fe[7] >> 25
	fe[8] += carry[7]
	fe[7] -= carry[7] << 25
	carry[8] = fe[8] >> 26
	fe[9] += carry[8]
	fe[8] -= carry[8] << 26
	carry[9] = fe[9] >> 25
	fe[9] -= carry[9] << 25

	s[0] = byte(fe[0] >> 0)
	s[1] = byte(fe[0] >> 8)
	s[2] = byte(fe[0] >> 16)
	s[3] = byte((fe[0] >> 24) | (fe[1] << 2))
	s[4] = byte(fe[1] >> 6)
	s[5] = byte(fe[1] >> 14)
	s[6] = byte((fe[1] >> 22) | (fe[2] << 3))
	s[7] = byte(fe[2] >> 5)
	s[8] = byte(fe[2] >> 13)
	s[9] = byte((fe[2] >> 21) | (fe[3] << 5))
	s[10] = byte(fe[3] >> 3)
	s[11] = byte(fe[3] >> 11)
	s[12] = byte((fe[3] >> 19) | (fe[4] << 6))
	s[13] = byte(fe[4] >> 2)
	s[14] = byte(fe[4] >> 10)
	s[15] = byte(fe[4] >> 18)
	s[16] = byte(fe[5] >> 0)
	s[17] = byte(fe[5] >> 8)
	s[18] = byte(fe[5] >> 16)
	s[19] = byte((fe[5] >> 24) | (fe[6] << 1))
	s[20] = byte(fe[6] >> 7)
	s[21] = byte(fe[6] >> 15)
	s[22] = byte((fe[6] >> 23) | (fe[7] << 3))
	s[23] = byte(fe[7] >> 5)
	s[24] = byte(fe[7] >> 13)
	s[25] = byte((fe[7] >> 21) | (fe[8] << 4))
	s[26] = byte(fe[8] >> 4)
	s[27] = byte(fe[8] >> 12)
	s[28] = byte((fe[8] >> 20) | (fe[9] << 6))
	s[29] = byte(fe[9] >> 2)
	s[30] = byte(fe[9] >> 10)
	s[31] = byte(fe[9] >> 18)
	return fe
}

// Sets fe to the little endian number encoded in buf modulo 2^255-19.
// Ignores the highest bit in buf.  Returns fe.
func (fe *FieldElement) SetBytes(buf *[32]byte) *FieldElement {
	return fe.setReduced(
		load4(buf[:]),
		load3(buf[4:])<<6,
		load3(buf[7:])<<5,
		load3(buf[10:])<<3,
		load3(buf[13:])<<2,
		load4(buf[16:]),
		load3(buf[20:])<<7,
		load3(buf[23:])<<5,
		load3(buf[26:])<<4,
		(load3(buf[29:])&8388607)<<2,
	)
}

// Sets fe to a * b.  Returns fe.
func (fe *FieldElement) Mul(a, b *FieldElement) *FieldElement {
	a0 := int64(a[0])
	a1 := int64(a[1])
	a2 := int64(a[2])
	a3 := int64(a[3])
	a4 := int64(a[4])
	a5 := int64(a[5])
	a6 := int64(a[6])
	a7 := int64(a[7])
	a8 := int64(a[8])
	a9 := int64(a[9])

	a1_2 := int64(2 * a[1])
	a3_2 := int64(2 * a[3])
	a5_2 := int64(2 * a[5])
	a7_2 := int64(2 * a[7])
	a9_2 := int64(2 * a[9])

	b0 := int64(b[0])
	b1 := int64(b[1])
	b2 := int64(b[2])
	b3 := int64(b[3])
	b4 := int64(b[4])
	b5 := int64(b[5])
	b6 := int64(b[6])
	b7 := int64(b[7])
	b8 := int64(b[8])
	b9 := int64(b[9])

	b1_19 := int64(19 * b[1])
	b2_19 := int64(19 * b[2])
	b3_19 := int64(19 * b[3])
	b4_19 := int64(19 * b[4])
	b5_19 := int64(19 * b[5])
	b6_19 := int64(19 * b[6])
	b7_19 := int64(19 * b[7])
	b8_19 := int64(19 * b[8])
	b9_19 := int64(19 * b[9])

	h0 := a0*b0 + a1_2*b9_19 + a2*b8_19 + a3_2*b7_19 + a4*b6_19 + a5_2*b5_19 + a6*b4_19 + a7_2*b3_19 + a8*b2_19 + a9_2*b1_19
	h1 := a0*b1 + a1*b0 + a2*b9_19 + a3*b8_19 + a4*b7_19 + a5*b6_19 + a6*b5_19 + a7*b4_19 + a8*b3_19 + a9*b2_19
	h2 := a0*b2 + a1_2*b1 + a2*b0 + a3_2*b9_19 + a4*b8_19 + a5_2*b7_19 + a6*b6_19 + a7_2*b5_19 + a8*b4_19 + a9_2*b3_19
	h3 := a0*b3 + a1*b2 + a2*b1 + a3*b0 + a4*b9_19 + a5*b8_19 + a6*b7_19 + a7*b6_19 + a8*b5_19 + a9*b4_19
	h4 := a0*b4 + a1_2*b3 + a2*b2 + a3_2*b1 + a4*b0 + a5_2*b9_19 + a6*b8_19 + a7_2*b7_19 + a8*b6_19 + a9_2*b5_19
	h5 := a0*b5 + a1*b4 + a2*b3 + a3*b2 + a4*b1 + a5*b0 + a6*b9_19 + a7*b8_19 + a8*b7_19 + a9*b6_19
	h6 := a0*b6 + a1_2*b5 + a2*b4 + a3_2*b3 + a4*b2 + a5_2*b1 + a6*b0 + a7_2*b9_19 + a8*b8_19 + a9_2*b7_19
	h7 := a0*b7 + a1*b6 + a2*b5 + a3*b4 + a4*b3 + a5*b2 + a6*b1 + a7*b0 + a8*b9_19 + a9*b8_19
	h8 := a0*b8 + a1_2*b7 + a2*b6 + a3_2*b5 + a4*b4 + a5_2*b3 + a6*b2 + a7_2*b1 + a8*b0 + a9_2*b9_19
	h9 := a0*b9 + a1*b8 + a2*b7 + a3*b6 + a4*b5 + a5*b4 + a6*b3 + a7*b2 + a8*b1 + a9*b0

	return fe.setReduced(h0, h1, h2, h3, h4, h5, h6, h7, h8, h9)
}

// Returns the unnormalized coefficients of fe^2.
func (fe *FieldElement) square() (h0, h1, h2, h3, h4, h5, h6, h7, h8, h9 int64) {
	f0 := int64(fe[0])
	f1 := int64(fe[1])
	f2 := int64(fe[2])
	f3 := int64(fe[3])
	f4 := int64(fe[4])
	f5 := int64(fe[5])
	f6 := int64(fe[6])
	f7 := int64(fe[7])
	f8 := int64(fe[8])
	f9 := int64(fe[9])
	f0_2 := int64(2 * fe[0])
	f1_2 := int64(2 * fe[1])
	f2_2 := int64(2 * fe[2])
	f3_2 := int64(2 * fe[3])
	f4_2 := int64(2 * fe[4])
	f5_2 := int64(2 * fe[5])
	f6_2 := int64(2 * fe[6])
	f7_2 := int64(2 * fe[7])
	f5_38 := 38 * f5
	f6_19 := 19 * f6
	f7_38 := 38 * f7
	f8_19 := 19 * f8
	f9_38 := 38 * f9

	h0 = f0*f0 + f1_2*f9_38 + f2_2*f8_19 + f3_2*f7_38 + f4_2*f6_19 + f5*f5_38
	h1 = f0_2*f1 + f2*f9_38 + f3_2*f8_19 + f4*f7_38 + f5_2*f6_19
	h2 = f0_2*f2 + f1_2*f1 + f3_2*f9_38 + f4_2*f8_19 + f5_2*f7_38 + f6*f6_19
	h3 = f0_2*f3 + f1_2*f2 + f4*f9_38 + f5_2*f8_19 + f6*f7_38
	h4 = f0_2*f4 + f1_2*f3_2 + f2*f2 + f5_2*f9_38 + f6_2*f8_19 + f7*f7_38
	h5 = f0_2*f5 + f1_2*f4 + f2_2*f3 + f6*f9_38 + f7_2*f8_19
	h6 = f0_2*f6 + f1_2*f5_2 + f2_2*f4 + f3_2*f3 + f7_2*f9_38 + f8*f8_19
	h7 = f0_2*f7 + f1_2*f6 + f2_2*f5 + f3_2*f4 + f8*f9_38
	h8 = f0_2*f8 + f1_2*f7_2 + f2_2*f6 + f3_2*f5_2 + f4*f4 + f9*f9_38
	h9 = f0_2*f9 + f1_2*f8 + f2_2*f7 + f3_2*f6 + f4_2*f5

	return
}

// Sets fe to a^2.  Returns fe.
func (fe *FieldElement) Square(a *FieldElement) *FieldElement {
	h0, h1, h2, h3, h4, h5, h6, h7, h8, h9 := a.square()
	return fe.setReduced(h0, h1, h2, h3, h4, h5, h6, h7, h8, h9)
}

// Sets fe to 2 * a^2.  Returns fe.
func (fe *FieldElement) DoubledSquare(a *FieldElement) *FieldElement {
	h0, h1, h2, h3, h4, h5, h6, h7, h8, h9 := a.square()
	h0 += h0
	h1 += h1
	h2 += h2
	h3 += h3
	h4 += h4
	h5 += h5
	h6 += h6
	h7 += h7
	h8 += h8
	h9 += h9
	return fe.setReduced(h0, h1, h2, h3, h4, h5, h6, h7, h8, h9)
}

// Set fe to the inverse of a.  Return fe.
func (fe *FieldElement) Inverse(a *FieldElement) *FieldElement {
	var t0, t1, t2, t3 FieldElement
	var i int

	t0.Square(a)
	t1.Square(&t0)
	for i = 1; i < 2; i++ {
		t1.Square(&t1)
	}
	t1.Mul(a, &t1)
	t0.Mul(&t0, &t1)
	t2.Square(&t0)
	t1.Mul(&t1, &t2)
	t2.Square(&t1)
	for i = 1; i < 5; i++ {
		t2.Square(&t2)
	}
	t1.Mul(&t2, &t1)
	t2.Square(&t1)
	for i = 1; i < 10; i++ {
		t2.Square(&t2)
	}
	t2.Mul(&t2, &t1)
	t3.Square(&t2)
	for i = 1; i < 20; i++ {
		t3.Square(&t3)
	}
	t2.Mul(&t3, &t2)
	t2.Square(&t2)
	for i = 1; i < 10; i++ {
		t2.Square(&t2)
	}
	t1.Mul(&t2, &t1)
	t2.Square(&t1)
	for i = 1; i < 50; i++ {
		t2.Square(&t2)
	}
	t2.Mul(&t2, &t1)
	t3.Square(&t2)
	for i = 1; i < 100; i++ {
		t3.Square(&t3)
	}
	t2.Mul(&t3, &t2)
	t2.Square(&t2)
	for i = 1; i < 50; i++ {
		t2.Square(&t2)
	}
	t1.Mul(&t2, &t1)
	t1.Square(&t1)
	for i = 1; i < 5; i++ {
		t1.Square(&t1)
	}
	return fe.Mul(&t1, &t0)
}

// Set fe to -x if x is negative and x otherwise.  Returns fe.
func (fe *FieldElement) Abs(x *FieldElement) *FieldElement {
	var xNeg FieldElement
	xNeg.Neg(x)
	fe.Set(x)
	fe.ConditionalSet(&xNeg, x.IsNegativeI())
	return fe
}

// Returns 1 if fe is negative, otherwise 0.
func (fe *FieldElement) IsNegativeI() int32 {
	var buf [32]byte
	fe.BytesInto(&buf)
	return int32(buf[0] & 1)
}

// Returns 1 if fe is non-zero, otherwise 0.
func (fe *FieldElement) IsNonZeroI() int32 {
	var ret uint8
	var buf [32]byte
	fe.BytesInto(&buf)
	for _, b := range buf {
		ret |= b
	}
	ret |= ret >> 4
	ret |= ret >> 2
	ret |= ret >> 1
	return int32(ret & 1)
}

// Returns 1 if fe is equal to one, otherwise 0.
func (fe *FieldElement) IsOneI() int32 {
	var b FieldElement
	return 1 - b.sub(fe, &feOne).IsNonZeroI()
}

// Returns 1 if fe is equal to a, otherwise 0.
func (fe *FieldElement) EqualsI(a *FieldElement) int32 {
	var b FieldElement
	return 1 - b.sub(fe, a).IsNonZeroI()
}

// Returns whether fe equals a.
func (fe *FieldElement) Equals(a *FieldElement) bool {
	var b FieldElement
	return b.sub(fe, a).IsNonZeroI() == 0
}

// Returns fe as a big.Int.
//
// WARNING Operations on big.Ints are not constant-time: do not use them
//         for cryptography unless you're sure this is not an issue.
func (fe *FieldElement) BigInt() *big.Int {
	var ret big.Int
	var buf, rBuf [32]byte
	fe.BytesInto(&buf)
	for i := 0; i < 32; i++ {
		rBuf[i] = buf[31-i]
	}
	return ret.SetBytes(rBuf[:])
}

// Sets fe to x modulo 2^255-19.
//
// WARNING Operations on big.Ints are not constant-time: do not use them
//         for cryptography unless you're sure this is not an issue.
func (fe *FieldElement) SetBigInt(x *big.Int) *FieldElement {
	var v, bi25519 big.Int
	bi25519.SetString(
		"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed", 16)
	buf := v.Mod(x, &bi25519).Bytes()
	var rBuf [32]byte
	for i := 0; i < len(buf) && i < 32; i++ {
		rBuf[i] = buf[len(buf)-i-1]
	}
	return fe.SetBytes(&rBuf)
}

// Sets fe to a^( ((2^255 - 19) - 5) / 8 ) = a^(2^252 - 3).  Returns fe.
//
// This method is useful to compute (inverse) square-roots
// with the method of Lagrange.
func (fe *FieldElement) Exp22523(a *FieldElement) *FieldElement {
	var t0, t1, t2 FieldElement
	var i int

	t0.Square(a)
	for i = 1; i < 1; i++ {
		t0.Square(&t0)
	}
	t1.Square(&t0)
	for i = 1; i < 2; i++ {
		t1.Square(&t1)
	}
	t1.Mul(a, &t1)
	t0.Mul(&t0, &t1)
	t0.Square(&t0)
	for i = 1; i < 1; i++ {
		t0.Square(&t0)
	}
	t0.Mul(&t1, &t0)
	t1.Square(&t0)
	for i = 1; i < 5; i++ {
		t1.Square(&t1)
	}
	t0.Mul(&t1, &t0)
	t1.Square(&t0)
	for i = 1; i < 10; i++ {
		t1.Square(&t1)
	}
	t1.Mul(&t1, &t0)
	t2.Square(&t1)
	for i = 1; i < 20; i++ {
		t2.Square(&t2)
	}
	t1.Mul(&t2, &t1)
	t1.Square(&t1)
	for i = 1; i < 10; i++ {
		t1.Square(&t1)
	}
	t0.Mul(&t1, &t0)
	t1.Square(&t0)
	for i = 1; i < 50; i++ {
		t1.Square(&t1)
	}
	t1.Mul(&t1, &t0)
	t2.Square(&t1)
	for i = 1; i < 100; i++ {
		t2.Square(&t2)
	}
	t1.Mul(&t2, &t1)
	t1.Square(&t1)
	for i = 1; i < 50; i++ {
		t1.Square(&t1)
	}
	t0.Mul(&t1, &t0)
	t0.Square(&t0)
	for i = 1; i < 2; i++ {
		t0.Square(&t0)
	}
	return fe.Mul(&t0, a)
}

// Sets fe to 1/sqrt(a).  Requires a to be a square.  Returns fe.
func (fe *FieldElement) InvSqrt(a *FieldElement) *FieldElement {
	var den2, den3, den4, den6, chk, t, t2 FieldElement
	den2.Square(a)
	den3.Mul(&den2, a)
	den4.Square(&den2)
	den6.Mul(&den2, &den4)
	t.Mul(&den6, a)
	t.Exp22523(&t)
	t.Mul(&t, &den3)
	t2.Mul(&t, &feI)

	chk.Square(&t)
	chk.Mul(&chk, a)

	fe.Set(&t)
	fe.ConditionalSet(&t2, 1-chk.IsOneI())
	return fe
}

// Sets fe to either 1/sqrt(a) or 1/sqrt(i*a).  Returns 1 in the former case
// and 0 in the latter.
func (fe *FieldElement) InvSqrtI(a *FieldElement) int32 {
	var inCaseA, inCaseB, inCaseD int32
	var den2, den3, den4, den6, chk, t, corr FieldElement
	den2.Square(a)
	den3.Mul(&den2, a)
	den4.Square(&den2)
	den6.Mul(&den2, &den4)
	t.Mul(&den6, a)
	t.Exp22523(&t)
	t.Mul(&t, &den3)

	// case       A           B            C             D
	// ---------------------------------------------------------------
	// t          1/sqrt(a)   -i/sqrt(a)   1/sqrt(i*a)   -i/sqrt(i*a)
	// chk        1           -1           -i            i
	// corr       1           i            1             i
	// ret        1           1            0             0

	chk.Square(&t)
	chk.Mul(&chk, a)

	inCaseA = chk.IsOneI()
	inCaseD = chk.EqualsI(&feI)
	chk.Neg(&chk)
	inCaseB = chk.IsOneI()

	corr.SetOne()
	corr.ConditionalSet(&feI, inCaseB+inCaseD)
	t.Mul(&t, &corr)
	fe.Set(&t)

	return inCaseA + inCaseB
}
