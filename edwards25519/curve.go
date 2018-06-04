// Go implementation of the elliptic curve Edwards25519 of which the
// Ristretto group is a subquotient.
package edwards25519

var (
	// parameter d of Edwards25519
	feD = FieldElement{
		-10913610, 13857413, -15372611, 6949391, 114729,
		-8787816, -6275908, -3247719, -18696448, -12055116,
	}

	// double feD
	fe2D = FieldElement{
		-21827239, -5839606, -30745221, 13898782, 229458,
		15978800, -12551817, -6495438, 29715968, 9444199,
	}

	// 1 / sqrt(-1-d)
	feInvSqrtMinusDMinusOne = FieldElement{
		-6111485, -4156064, 27798727, -12243468, 25904040,
		-120897, -20826367, 7060776, -6093568, 1986012,
	}

	// (d-1)^2
	feDMinusOneSquared = FieldElement{
		15551795, -11097455, -13425098, -10125071, -11896535,
		10178284, -26634327, 4729244, -5282110, -10116402,
	}

	// 1 - d^2
	feOneMinusDSquared = FieldElement{
		6275446, -16617371, -22938544, -3773710, 11667077,
		7397348, -27922721, 1766195, -24433858, 672203,
	}

	epZero = ExtendedPoint{feZero, feOne, feOne, feZero}
)

// // (X:Y:Z) satisfying x=X/Z, y=Y/Z.  Aka P2.
// type ProjectivePoint struct {
// 	X, Y, Z FieldElement
// }

// (X:Y:Z:T) satisfying x=X/Z, y=Y/Z, X*Y=Z*T.  Aka P3.
type ExtendedPoint struct {
	X, Y, Z, T FieldElement
}

// ((X:Z),(Y:T)) satisfying x=X/Z, y=Y/T. Aka P1P1.
type CompletedPoint struct {
	X, Y, Z, T FieldElement
}

// Set p to zero, the neutral element.  Return p.
func (p *ExtendedPoint) SetZero() *ExtendedPoint {
	p.X.SetZero()
	p.Y.SetOne()
	p.Z.SetOne()
	p.T.SetZero()
	return p
}

// Set p to q.  Returns p.
func (p *ExtendedPoint) Set(q *ExtendedPoint) *ExtendedPoint {
	p.X.Set(&q.X)
	p.Y.Set(&q.Y)
	p.Z.Set(&q.Z)
	p.T.Set(&q.T)
	return p
}

// Set p to q if b == 1.  Assumes b is 0 or 1.   Returns p.
func (p *ExtendedPoint) ConditionalSet(q *ExtendedPoint, b int32) *ExtendedPoint {
	p.X.ConditionalSet(&q.X, b)
	p.Y.ConditionalSet(&q.Y, b)
	p.Z.ConditionalSet(&q.Z, b)
	p.T.ConditionalSet(&q.T, b)
	return p
}

// Set p to the point corresponding to the given point (s,t) on the
// associated Jacobi quartic.
func (p *CompletedPoint) SetJacobiQuartic(s, t *FieldElement) *CompletedPoint {
	var s2 FieldElement
	s2.Square(s)

	// Set x to 2 * s * 1/sqrt(-d-1)
	p.X.double(s)
	p.X.Mul(&p.X, &feInvSqrtMinusDMinusOne)

	// Set z to t
	p.Z.Set(t)

	// Set y to 1-s^2
	p.Y.sub(&feOne, &s2)

	// Set t to 1+s^2
	p.T.add(&feOne, &s2)
	return p
}

// Set p to the curvepoint corresponding to r0 via Mike Hamburg's variation
// on Elligator2 for Ristretto.  Returns p.
func (p *CompletedPoint) SetRistrettoElligator2(r0 *FieldElement) *CompletedPoint {
	var r, rPlusD, rPlusOne, D, N, ND, sqrt, twiddle, sgn FieldElement
	var s, t, rSubOne, r0i, sNeg FieldElement

	var b int32

	// r := i * r0^2
	r0i.Mul(r0, &feI)
	r.Mul(r0, &r0i)

	// D := -((d*r)+1) * (r + d)
	rPlusD.add(&feD, &r)
	D.Mul(&feD, &r)
	D.add(&D, &feOne)
	D.Mul(&D, &rPlusD)
	D.Neg(&D)

	// N := -(d^2 - 1)(r + 1)
	rPlusOne.add(&r, &feOne)
	N.Mul(&feOneMinusDSquared, &rPlusOne)

	// sqrt is the inverse square root of N*D or of i*N*D.
	// b=1 iff n1 is square.
	ND.Mul(&N, &D)
	b = sqrt.InvSqrtI(&ND)
	sqrt.Abs(&sqrt)

	twiddle.SetOne()
	twiddle.ConditionalSet(&r0i, 1-b)
	sgn.SetOne()
	sgn.ConditionalSet(&feMinusOne, 1-b)
	sqrt.Mul(&sqrt, &twiddle)

	// s = N * sqrt(N*D) * twiddle
	s.Mul(&sqrt, &N)

	// t = -sgn * sqrt * s * (r-1) * (d-1)^2 - 1
	t.Neg(&sgn)
	t.Mul(&sqrt, &t)
	t.Mul(&s, &t)
	t.Mul(&feDMinusOneSquared, &t)
	rSubOne.sub(&r, &feOne)
	t.Mul(&rSubOne, &t)
	t.sub(&t, &feOne)

	sNeg.Neg(&s)
	s.ConditionalSet(&sNeg, equal30(s.IsNegativeI(), b))
	return p.SetJacobiQuartic(&s, &t)
}

// Sets p to q+r.  Returns p
func (p *CompletedPoint) AddExtended(q, r *ExtendedPoint) *CompletedPoint {
	var a, b, c, d, t FieldElement

	a.sub(&q.Y, &q.X)
	t.sub(&r.Y, &r.X)
	a.Mul(&a, &t)
	b.add(&q.X, &q.Y)
	t.add(&r.X, &r.Y)
	b.Mul(&b, &t)
	c.Mul(&q.T, &r.T)
	c.Mul(&c, &fe2D)
	d.Mul(&q.Z, &r.Z)
	d.add(&d, &d)
	p.X.sub(&b, &a)
	p.T.sub(&d, &c)
	p.Z.add(&d, &c)
	p.Y.add(&b, &a)

	return p
}

// Set p to 2 * q.  Returns p.
func (p *CompletedPoint) DoubleExtended(q *ExtendedPoint) *CompletedPoint {
	var a, b, c, d FieldElement

	a.Square(&q.X)
	b.Square(&q.Y)
	c.DoubledSquare(&q.Z)
	d.Neg(&a)
	p.X.add(&q.X, &q.Y)
	p.X.Square(&p.X)
	p.X.sub(&p.X, &a)
	p.X.sub(&p.X, &b)
	p.Z.add(&d, &b)
	p.T.sub(&p.Z, &c)
	p.Y.sub(&d, &b)

	return p
}

// Set p to 2 * q. Returns p.
func (p *ExtendedPoint) Double(q *ExtendedPoint) *ExtendedPoint {
	var tmp CompletedPoint
	tmp.DoubleExtended(q)
	p.SetCompleted(&tmp)
	return p
}

// Set p to q + r. Returns p.
func (p *ExtendedPoint) Add(q, r *ExtendedPoint) *ExtendedPoint {
	var tmp CompletedPoint
	tmp.AddExtended(q, r)
	p.SetCompleted(&tmp)
	return p
}

// Sets p to q.  Returns p.
func (p *ExtendedPoint) SetCompleted(q *CompletedPoint) *ExtendedPoint {
	p.X.Mul(&q.X, &q.T)
	p.Y.Mul(&q.Y, &q.Z)
	p.Z.Mul(&q.Z, &q.T)
	p.T.Mul(&q.X, &q.Y)
	return p
}

// Set p to a point corresponding to the encoded group element of
// the ristretto group.  Returns whether the buffer encoded a group element.
func (p *ExtendedPoint) SetRistretto(buf *[32]byte) bool {
	var s, s2, chk, yDen, yNum, yDen2, xDen2, isr, xDenInv FieldElement
	var yDenInv, t FieldElement
	var b, ret int32

	s.SetBytes(buf)
	ret = s.IsNegativeI()
	s2.Square(&s)
	yDen.add(&feOne, &s2)
	yNum.sub(&feOne, &s2)
	yDen2.Square(&yDen)
	xDen2.Square(&yNum)
	xDen2.Mul(&xDen2, &feD)
	xDen2.add(&xDen2, &yDen2)
	xDen2.Neg(&xDen2)
	t.Mul(&xDen2, &yDen2)
	isr.InvSqrt(&t)
	chk.Square(&isr)
	chk.Mul(&chk, &t)
	ret |= 1 - chk.IsOneI()
	xDenInv.Mul(&isr, &yDen)
	yDenInv.Mul(&xDenInv, &isr)
	yDenInv.Mul(&yDenInv, &xDen2)
	p.X.Mul(&s, &xDenInv)
	p.X.add(&p.X, &p.X)
	b = p.X.IsNegativeI()
	t.Neg(&p.X)
	p.X.ConditionalSet(&t, b)
	p.Y.Mul(&yNum, &yDenInv)
	p.Z.SetOne()
	p.T.Mul(&p.X, &p.Y)
	ret |= p.T.IsNegativeI()
	ret |= 1 - p.Y.IsNonZeroI()
	p.X.ConditionalSet(&feZero, ret)
	p.Y.ConditionalSet(&feZero, ret)
	p.Z.ConditionalSet(&feZero, ret)
	p.T.ConditionalSet(&feZero, ret)
	return ret == 0
}

// Pack p using the Ristretto encoding and return it.
// Requires p to be even.
func (p *ExtendedPoint) Ristretto() []byte {
	var buf [32]byte
	p.RistrettoInto(&buf)
	return buf[:]
}

// Pack p using the Ristretto encoding and write to buf.  Returns p.
// Requires p to be even.
func (p *ExtendedPoint) RistrettoInto(buf *[32]byte) *ExtendedPoint {
	var d, u1, u2, isr, i1, i2, zInv, denInv, nx, ny, s FieldElement
	var b int32

	d.add(&p.Z, &p.Y)
	u1.sub(&p.Z, &p.Y)
	u1.Mul(&u1, &d)

	u2.Mul(&p.X, &p.Y)

	isr.Square(&u2)
	isr.Mul(&isr, &u1)
	isr.InvSqrt(&isr)

	i1.Mul(&isr, &u1)
	i2.Mul(&isr, &u2)

	zInv.Mul(&i1, &i2)
	zInv.Mul(&zInv, &p.T)

	d.Mul(&zInv, &p.T)

	nx.Mul(&p.Y, &feI)
	ny.Mul(&p.X, &feI)
	denInv.Mul(&feInvSqrtMinusDMinusOne, &i1)

	b = 1 - d.IsNegativeI()
	nx.ConditionalSet(&p.X, b)
	ny.ConditionalSet(&p.Y, b)
	denInv.ConditionalSet(&i2, b)

	d.Mul(&nx, &zInv)
	b = d.IsNegativeI()
	d.Neg(&ny)
	ny.ConditionalSet(&d, b)

	s.sub(&p.Z, &ny)
	s.Mul(&s, &denInv)

	b = s.IsNegativeI()
	d.Neg(&s)
	s.ConditionalSet(&d, b)

	s.BytesInto(buf)
	return p
}

// Compute 5-bit window for the scalar s.
func computeScalarWindow5(s *[32]byte, w *[51]int8) {
	for i := 0; i < 6; i++ {
		w[8*i+0] = int8(s[5*i+0] & 31)
		w[8*i+1] = int8((s[5*i+0] >> 5) & 31)
		w[8*i+1] ^= int8((s[5*i+1] << 3) & 31)
		w[8*i+2] = int8((s[5*i+1] >> 2) & 31)
		w[8*i+3] = int8((s[5*i+1] >> 7) & 31)
		w[8*i+3] ^= int8((s[5*i+2] << 1) & 31)
		w[8*i+4] = int8((s[5*i+2] >> 4) & 31)
		w[8*i+4] ^= int8((s[5*i+3] << 4) & 31)
		w[8*i+5] = int8((s[5*i+3] >> 1) & 31)
		w[8*i+6] = int8((s[5*i+3] >> 6) & 31)
		w[8*i+6] ^= int8((s[5*i+4] << 2) & 31)
		w[8*i+7] = int8((s[5*i+4] >> 3) & 31)
	}
	w[8*6+0] = int8(s[5*6+0] & 31)
	w[8*6+1] = int8((s[5*6+0] >> 5) & 31)
	w[8*6+1] ^= int8((s[5*6+1] << 3) & 31)
	w[8*6+2] = int8((s[5*6+1] >> 2) & 31)

	/* Making it signed */
	var carry int8 = 0
	for i := 0; i < 50; i++ {
		w[i] += carry
		w[i+1] += w[i] >> 5
		w[i] &= 31
		carry = w[i] >> 4
		w[i] -= carry << 5
	}
	w[50] += carry
}

// Set p to s * q.  Returns p.
func (p *ExtendedPoint) ScalarMult(q *ExtendedPoint, s *[32]byte) *ExtendedPoint {
	// See eg. https://cryptojedi.org/peter/data/eccss-20130911b.pdf
	var lut [17]ExtendedPoint
	var t ExtendedPoint
	var window [51]int8

	// Precomputations.  TODO expose an API for such precomputations
	computeScalarWindow5(s, &window)
	lut[0].SetZero()
	lut[1].Set(q)
	for i := 2; i < 16; i += 2 {
		lut[i].Double(&lut[i>>1])
		lut[i+1].Add(&lut[i], q)
	}
	lut[16].Double(&lut[8])

	// Compute!
	p.SetZero()
	for i := 50; i >= 0; i-- {
		for j := 0; j < 5; j++ {
			p.Double(p)
		}

		t.Set(&lut[0])
		b := int32(window[i])
		for j := 1; j <= 16; j++ {
			c := equal15(b, int32(-j)) | equal15(b, int32(j))
			t.ConditionalSet(&lut[j], c)
		}
		var v FieldElement
		c := negative(b)
		v.Neg(&t.X)
		t.X.ConditionalSet(&v, c)
		v.Neg(&t.T)
		t.T.ConditionalSet(&v, c)

		p.Add(p, &t)
	}

	return p
}

// Sets p to -q.  Returns p.
func (p *ExtendedPoint) Neg(q *ExtendedPoint) *ExtendedPoint {
	p.X.Neg(&q.X)
	p.Y.Set(&q.Y)
	p.Z.Set(&q.Z)
	p.T.Neg(&q.T)
	return p
}
