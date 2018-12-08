package edwards25519

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestComputeScalar5NAF(t *testing.T) {
	rnd := rand.New(rand.NewSource(37))
	var rs, s [32]byte
	var sbi, wbi big.Int
	for i := 0; i < 1000; i++ {
		rnd.Read(s[0:32])
		s[31] &= 31
		var w [256]int8
		computeScalar5NAF(&s, &w)
		for j := 0; j < 32; j++ {
			rs[j] = s[31-j]
		}
		sbi.SetBytes(rs[:])
		var power, summand big.Int
		power.SetUint64(1)
		wbi.SetUint64(0)
		for j := 0; j < 255; j++ {
			summand.SetInt64(int64(w[j]))
			summand.Mul(&summand, &power)
			wbi.Add(&wbi, &summand)
			power.Add(&power, &power)
		}
		if wbi.Cmp(&sbi) != 0 {
			t.Fatalf("5NAF(%v) = %v  %v != %v", s, w, &sbi, &wbi)
		}
	}
}

func TestVarTimeScalarMult(t *testing.T) {
	rnd := rand.New(rand.NewSource(37))
	var fe FieldElement
	var cp CompletedPoint
	var q, p1, p2 ExtendedPoint
	var s [32]byte
	for i := 0; i < 1000; i++ {
		var buf [32]byte
		rnd.Read(buf[:])
		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		q.SetCompleted(&cp)
		rnd.Read(s[0:32])
		s[31] &= 31
		p1.ScalarMult(&q, &s)
		p2.VarTimeScalarMult(&q, &s)
		if p1.RistrettoEqualsI(&p2) != 1 {
			t.Fatalf("[%v]%v = %v != %v", s, q, p1, p2)
		}
	}
}

func BenchmarkVarTimeScalarMult(b *testing.B) {
	rnd := rand.New(rand.NewSource(37))
	var buf, s [32]byte
	var cp CompletedPoint
	var ep ExtendedPoint
	var fe FieldElement
	rnd.Read(s[0:32])
	s[31] &= 31
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ep.VarTimeScalarMult(&ep, &s)
	}
}
