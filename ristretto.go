// Pure Go implementation of the Ristretto prime-order group built from
// the Edwards curve Edwards25519.
//
// Many cryptographic schemes need a group of prime order.  Popular and
// efficient elliptic curves like (Edwards25519 of `ed25519` fame) are
// rarely of prime order.  There is, however, a convenient method
// to construct a prime order group from such curves, using a method
// called Ristretto proposed by Mike Hamburg.
package ristretto

import (
	"crypto/rand"

	"github.com/bwesterb/go-ristretto/edwards25519"
)

// Represents an element of the Ristretto group over Edwards25519.
type Point edwards25519.ExtendedPoint

// Sets p to zero (the neutral element).  Returns p.
func (p *Point) SetZero() *Point {
	p.e().SetZero()
	return p
}

// Sets p to q + r.  Returns p.
func (p *Point) Add(q, r *Point) *Point {
	p.e().Add(q.e(), r.e())
	return p
}

// Sets p to q - r.  Returns p.
func (p *Point) Sub(q, r *Point) *Point {
	// TODO optimize
	var negR Point
	negR.Neg(r)
	p.Add(q, &negR)
	return p
}

// Sets p to -q.  Returns p.
func (p *Point) Neg(q *Point) *Point {
	p.e().Neg(q.e())
	return p
}

// Packs p into the given buffer.  Returns p.
func (p *Point) BytesInto(buf *[32]byte) *Point {
	p.e().RistrettoInto(buf)
	return p
}

// Returns a packed version of p.
func (p *Point) Bytes() []byte {
	return p.e().Ristretto()
}

// Sets p to the point encoded in buf using Bytes().
// Not every input encodes a point.  Returns whether the buffer encoded a point.
func (p *Point) SetBytes(buf *[32]byte) bool {
	return p.e().SetRistretto(buf)
}

// Sets p to the point corresponding to buf using the Elligator2 encoding.
//
// In contrast to SetBytes():
//   1.  Every input buffer will decode to a point.
//   2.  SetElligator() is not injective: for every point there are
//       approximately four buffers that will encode to it.
func (p *Point) SetElligator(buf *[32]byte) *Point {
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	fe.SetBytes(buf)
	cp.SetRistrettoElligator2(&fe)
	p.e().SetCompleted(&cp)
	return p
}

// Sets p to s * q.  Returns p.
func (p *Point) ScalarMult(q *Point, s *Scalar) *Point {
	p.e().ScalarMult(q.e(), (*[32]uint8)(s))
	return p
}

// Sets p to a random point.  Returns p.
func (p *Point) Rand() *Point {
	var buf [32]byte
	rand.Read(buf[:])
	return p.SetElligator(&buf)
}

func (p *Point) e() *edwards25519.ExtendedPoint {
	return (*edwards25519.ExtendedPoint)(p)
}
