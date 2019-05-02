package edwards25519

import (
	"fmt"
)

// (S,T,Z) represents the point (S/Z,T/Z) on the associated Jacobi quartic.
type ProjectiveJacobiPoint struct {
	S, T, Z FieldElement
}

// Computes the at most 8 positive FieldElements f such that p == elligator2(f).
// Assumes p is even.
//
// Returns a bitmask of which elements in fes are set.
func (p *ExtendedPoint) RistrettoElligator2Inverse(fes *[8]FieldElement) uint8 {
	var setMask uint8
	var p2 ExtendedPoint
	var jc ProjectiveJacobiPoint

	// Elligator2 computes a Point from a FieldElement in two steps: first
	// it computes a (s,t) on the Jacobi quartic and then computes the
	// corresponding even point on the Edwards curve.
	//
	// We invert in three steps.  Any Ristretto point has four representatives
	// as even Edwards points.  For each of those even Edwards points,
	// there are two points on the Jacobi quartic that map to it.
	// Each of those eight points on the Jacobi quartic might have an
	// Elligator2 preimage.
	//
	// Essentially we first loop over the four representatives of our point,
	// then for each of them consider both points on the Jacobi quartic and
	// check whether they have an inverse under Elligator2.  We take a few
	// shortcuts though.
	//
	//  1. We only compute two Jacobi quartic points directly from the
	//     the representatives.  The other two can be derived from them.
	//  2. We reuse knowledge of positivity of s in the point (s,t) on
	//     the Jacobi Quartic for the dual point (-s,-t).

	for j := 0; j < 4; j++ {
		// We loop over the four even points in the same Ristretto equivalence
		// class as p = (x, y).
		//
		//    j == 0    p itself
		//    j == 1    (-x, -y)
		//    j == 2    (iy, ix)
		//    j == 3    (-iy, -ix)
		//
		// For each we compute the jacobi point.  For j == 0 and j == 2 we do
		// this directly from the representative.  For j == 1 we use the one
		// computed from j == 0 and similarly with j == 3.
		if j == 0 {
			p2.Set(p) // First one is p itself.
		} else if j == 2 {
			p2.X.Set(&p.Y)
			p2.Y.Set(&p.X)
			p2.Z.Mul(&p.Z, &feI)
			p2.T.Neg(&p.T)
		}

		if j == 0 || j == 2 {
			jc.SetExtended(&p2)
		} else { // j == 1 or j == 3
			jc2 := jc
			jc.S.Set(&jc2.Z)
			jc.T.Neg(&jc2.T)
			jc.Z.Set(&jc2.S)
		}

		ok := int(jc.Z.IsNonZeroI())

		// TODO reuse computation
		var s, zInv FieldElement
		zInv.Inverse(&jc.Z)
		s.Mul(&zInv, &jc.S)
		sPos := 1 - s.IsNegativeI()

		setMask |= uint8(jc.elligator2Inverse(&fes[2*j], sPos) & ok << uint(2*j))
		jc.Dual(&jc)
		setMask |= uint8(jc.elligator2Inverse(&fes[2*j+1], 1-sPos) & ok << uint(2*j+1))
	}
	return setMask
}

// Set p to the point correspoding to q on the associated Jacobi quartic.
// Returns p.
func (p *ProjectiveJacobiPoint) SetExtended(q *ExtendedPoint) *ProjectiveJacobiPoint {
	var Z2, Y2, ZmY, tmp FieldElement

	// TODO - use q.T
	//      - double-check X=0 cases

	// Z = X sqrt(Z^2 - Y^2)
	Z2.Square(&q.Z)
	Y2.Square(&q.Y)
	tmp.sub(&Z2, &Y2)
	tmp.Sqrt(&tmp)
	p.Z.Mul(&q.X, &tmp)

	// S = (Z-Y)X
	ZmY.sub(&q.Z, &q.Y)
	p.S.Mul(&ZmY, &q.X)

	// T = 2 Z q.Z (Z-Y) 1/sqrt(-d-1)
	tmp.double(&feInvSqrtMinusDMinusOne)
	tmp.Mul(&tmp, &q.Z)
	tmp.Mul(&tmp, &p.Z)
	p.T.Mul(&tmp, &ZmY)

	return p
}

func (p *ProjectiveJacobiPoint) Dual(q *ProjectiveJacobiPoint) *ProjectiveJacobiPoint {
	p.S.Neg(&q.S)
	p.T.Neg(&q.T)
	p.Z.Set(&q.Z)
	return p
}

// Elligator2 is defined in two steps: first a field element is converted
// to a point (s,t) on the Jacobi quartic associated to the Edwards curve.
// Then this point is mapped to a point on the Edwards curve.
// This function computes a field element that is mapped to a given (s,t)
// with Elligator2 if it exists.
//
// sPos should be 1 if s is positive and 0 if it is not.
// (A ProjectiveJacobiPoint doesn't store s directly, but rather Z and S
// with S = s Z and so it is expensive to check whether s is positive.)
//
// Returns 1 if a preimage is found and 0 if none exists.
func (p *ProjectiveJacobiPoint) elligator2Inverse(fe *FieldElement, sPos int32) int {
	var x, y, a, a2, S2, S4, Z2, invSqY, negS2 FieldElement

	ret := p.Z.IsNonZeroI()
	done := int32(0)

	Z2.Square(&p.Z)

	// Special case: S = 0.  If S is zero, either t = 1 or t = -1.
	// If t=1, then sqrt(i*d) is the preimage.  There is no preimage if t=-1.
	sNonZero := p.S.IsNonZeroI()
	tEqualsZ2 := p.T.EqualsI(&Z2) // T = Z^2 if and only if t = 1
	ret &= 1 - ((1 - sNonZero) & (1 - tEqualsZ2))
	fe.ConditionalSet(&feSqrtID, 1-sNonZero)
	done = 1 - sNonZero

	// a := (T + Z^2) (d+1)/(d-1) = (t+1) (d+1)/(d-1)
	a.add(&p.T, &Z2)
	a.Mul(&a, &feDp1OverDm1)
	a2.Square(&a)

	// y := 1/sqrt(i (S^4 - a^2)).
	S2.Square(&p.S)
	S4.Square(&S2)
	invSqY.sub(&S4, &a2)
	invSqY.Mul(&invSqY, &feI)

	sq := y.InvSqrtI(&invSqY)
	ret &= sq // there is no preimage if the square root does not exist
	done |= 1 - sq

	// x := (a + sign(s)*S^2) y
	negS2.Neg(&S2)
	S2.ConditionalSet(&negS2, 1-sPos)
	x.add(&a, &S2)
	x.Mul(&x, &y)

	// fe := abs(x)
	x.Abs(&x)
	fe.ConditionalSet(&x, 1-done)
	return int(ret)
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

	// s = N * sqrt * twiddle
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

// WARNING This operation is not constant-time.  Do not use for cryptography
//         unless you're sure this is not an issue.
func (p *ProjectiveJacobiPoint) String() string {
	return fmt.Sprintf("ProjectiveJacobiPoint(%v, %v, %v)", p.S, p.T, p.Z)
}
