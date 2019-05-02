package edwards25519_test

import (
	"math/big"
	"testing"

	"github.com/bwesterb/go-ristretto/edwards25519"
)

func TestAddExtendedNiels(t *testing.T) {
	var buf1, buf2 [32]byte
	var cp1, cp2, cp3 edwards25519.CompletedPoint
	var np2 edwards25519.NielsPoint
	var fe1, fe2 edwards25519.FieldElement
	var ep1, ep2, ep3a, ep3b edwards25519.ExtendedPoint
	for i := 0; i < 1000; i++ {
		rnd.Read(buf1[:])
		rnd.Read(buf2[:])
		fe1.SetBytes(&buf1)
		fe2.SetBytes(&buf2)
		cp1.SetRistrettoElligator2(&fe1)
		cp2.SetRistrettoElligator2(&fe2)
		ep1.SetCompleted(&cp1)
		ep2.SetCompleted(&cp2)
		ep3a.Add(&ep1, &ep2)
		np2.SetExtended(&ep2)
		cp3.AddExtendedNiels(&ep1, &np2)
		ep3b.SetCompleted(&cp3)
		if ep3a.RistrettoEqualsI(&ep3b) != 1 {
			t.Fatalf("%v + %v = %v != %v", ep1, ep2, ep3a, ep3b)
		}
	}
}

func TestSubExtendedNiels(t *testing.T) {
	var buf1, buf2 [32]byte
	var cp1, cp2, cp3 edwards25519.CompletedPoint
	var np2 edwards25519.NielsPoint
	var fe1, fe2 edwards25519.FieldElement
	var ep1, ep2, ep3a, ep3b edwards25519.ExtendedPoint
	for i := 0; i < 1000; i++ {
		rnd.Read(buf1[:])
		rnd.Read(buf2[:])
		fe1.SetBytes(&buf1)
		fe2.SetBytes(&buf2)
		cp1.SetRistrettoElligator2(&fe1)
		cp2.SetRistrettoElligator2(&fe2)
		ep1.SetCompleted(&cp1)
		ep2.SetCompleted(&cp2)
		ep3a.Sub(&ep1, &ep2)
		np2.SetExtended(&ep2)
		cp3.SubExtendedNiels(&ep1, &np2)
		ep3b.SetCompleted(&cp3)
		if ep3a.RistrettoEqualsI(&ep3b) != 1 {
			t.Fatalf("%v - %v = %v != %v", ep1, ep2, ep3a, ep3b)
		}
	}
}

func TestTableVarTimeBaseScalarMult(t *testing.T) {
	var table edwards25519.ScalarMultTable
	var B, p1, p2 edwards25519.ExtendedPoint
	B.SetBase()
	table.Compute(&B)
	var s [32]byte
	for i := 0; i < 1000; i++ {
		rnd.Read(s[:])
		s[31] &= 31
		table.ScalarMult(&p1, &s)
		table.VarTimeScalarMult(&p2, &s)
		if p1.RistrettoEqualsI(&p2) != 1 {
			t.Fatalf("[%v]B = %v != %v", s, p2, p1)
		}
	}
}

func TestTableBaseScalarMult(t *testing.T) {
	var table edwards25519.ScalarMultTable
	var B, p1, p2 edwards25519.ExtendedPoint
	B.SetBase()
	table.Compute(&B)
	var s [32]byte
	for i := 0; i < 1000; i++ {
		rnd.Read(s[:])
		s[31] &= 31
		table.ScalarMult(&p1, &s)
		p2.ScalarMult(&B, &s)
		if p1.RistrettoEqualsI(&p2) != 1 {
			t.Fatalf("[%v]B = %v != %v", s, p2, p1)
		}
	}
}

func TestBaseScalarMultTable(t *testing.T) {
	var table edwards25519.ScalarMultTable
	var B edwards25519.ExtendedPoint
	B.SetBase()
	table.Compute(&B)
	if edwards25519.BaseScalarMultTable != table {
		t.Fatalf("BaseScalarMultTable has incorrect values")
	}
}

func BenchmarkScalarMultTableCompute(b *testing.B) {
	var buf [32]byte
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var table edwards25519.ScalarMultTable
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		table.Compute(&ep)
	}
}

func BenchmarkScalarMultTableScalarMult(b *testing.B) {
	var buf, sBuf [32]byte
	var biS big.Int
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var fe edwards25519.FieldElement
	var table edwards25519.ScalarMultTable
	biS.Rand(rnd, &biL)
	srBuf := biS.Bytes()
	for j := 0; j < len(srBuf); j++ {
		sBuf[j] = srBuf[len(srBuf)-j-1]
	}
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	table.Compute(&ep)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		table.ScalarMult(&ep, &sBuf)
	}
}

func BenchmarkScalarMultTableVarTimeScalarMult(b *testing.B) {
	var buf, sBuf [32]byte
	var biS big.Int
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var fe edwards25519.FieldElement
	var table edwards25519.ScalarMultTable
	biS.Rand(rnd, &biL)
	srBuf := biS.Bytes()
	for j := 0; j < len(srBuf); j++ {
		sBuf[j] = srBuf[len(srBuf)-j-1]
	}
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	table.Compute(&ep)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		table.VarTimeScalarMult(&ep, &sBuf)
	}
}
