package edwards25519_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/bwesterb/go-ristretto/cref"
	"github.com/bwesterb/go-ristretto/edwards25519"
)

func TestElligatorAndRistretto(t *testing.T) {
	var buf, goBuf, cBuf, goBuf2, cBuf2 [32]byte
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var ep2 edwards25519.ExtendedPoint

	var cP cref.GroupGe
	var cP2 cref.GroupGe
	var cFe cref.Fe25519

	for i := 0; i < 1000; i++ {
		rnd.Read(buf[:])

		cFe.Unpack(&buf)
		cP.Elligator(&cFe)
		cP.Pack(&cBuf)

		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
		ep.RistrettoInto(&goBuf)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("pack o elligator ( %v ) = %v != %v", buf, cBuf, goBuf)
		}

		ep2.SetRistretto(&goBuf)
		ep2.RistrettoInto(&goBuf2)

		cP2.Unpack(&cBuf)
		cP2.Pack(&cBuf2)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("pack o unpack o pack o elligator ( %v ) = %v != %v",
				buf, cBuf2, goBuf2)
		}
	}
}

func TestPointDouble(t *testing.T) {
	var buf, cBuf, goBuf [32]byte
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep, ep2 edwards25519.ExtendedPoint

	var cFe cref.Fe25519
	var cP, cP2 cref.GroupGe

	for i := 0; i < 1000; i++ {
		rnd.Read(buf[:])

		cFe.Unpack(&buf)
		cP.Elligator(&cFe)
		cP2.Double(&cP)
		cP2.Pack(&cBuf)

		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
		ep2.Double(&ep)
		ep2.RistrettoInto(&goBuf)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("2*%v = %v != %v", ep, ep2, cP2)
		}
	}
}

func TestPointSub(t *testing.T) {
	var buf1, buf2, cBuf, goBuf [32]byte
	var fe1, fe2 edwards25519.FieldElement
	var cp1, cp2 edwards25519.CompletedPoint
	var ep1, ep2, ep3 edwards25519.ExtendedPoint

	var cFe1, cFe2 cref.Fe25519
	var cP1, cP2, cP3 cref.GroupGe

	for i := 0; i < 1000; i++ {
		rnd.Read(buf1[:])
		rnd.Read(buf2[:])

		cFe1.Unpack(&buf1)
		cFe2.Unpack(&buf2)
		cP1.Elligator(&cFe1)
		cP2.Elligator(&cFe2)
		cP2.Neg(&cP2)
		cP3.Add(&cP1, &cP2)
		cP3.Pack(&cBuf)

		fe1.SetBytes(&buf1)
		fe2.SetBytes(&buf2)
		cp1.SetRistrettoElligator2(&fe1)
		cp2.SetRistrettoElligator2(&fe2)
		ep1.SetCompleted(&cp1)
		ep2.SetCompleted(&cp2)
		ep3.Sub(&ep1, &ep2)
		ep3.RistrettoInto(&goBuf)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("%v - %v = %v != %v", ep1, ep2, ep3, cP3)
		}
	}
}

func TestPointAdd(t *testing.T) {
	var buf1, buf2, cBuf, goBuf [32]byte
	var fe1, fe2 edwards25519.FieldElement
	var cp1, cp2 edwards25519.CompletedPoint
	var ep1, ep2, ep3 edwards25519.ExtendedPoint

	var cFe1, cFe2 cref.Fe25519
	var cP1, cP2, cP3 cref.GroupGe

	for i := 0; i < 1000; i++ {
		rnd.Read(buf1[:])
		rnd.Read(buf2[:])

		cFe1.Unpack(&buf1)
		cFe2.Unpack(&buf2)
		cP1.Elligator(&cFe1)
		cP2.Elligator(&cFe2)
		cP3.Add(&cP1, &cP2)
		cP3.Pack(&cBuf)

		fe1.SetBytes(&buf1)
		fe2.SetBytes(&buf2)
		cp1.SetRistrettoElligator2(&fe1)
		cp2.SetRistrettoElligator2(&fe2)
		ep1.SetCompleted(&cp1)
		ep2.SetCompleted(&cp2)
		ep3.Add(&ep1, &ep2)
		ep3.RistrettoInto(&goBuf)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("%v + %v = %v != %v", ep1, ep2, ep3, cP3)
		}
	}
}

func TestScalarMult(t *testing.T) {
	var buf, sBuf, cBuf, goBuf [32]byte
	var biS big.Int
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep, ep2 edwards25519.ExtendedPoint

	var cFe cref.Fe25519
	var cS cref.GroupScalar
	var cP, cP2 cref.GroupGe

	for i := 0; i < 1000; i++ {
		rnd.Read(buf[:])
		biS.Rand(rnd, &biL)
		srBuf := biS.Bytes()
		for j := 0; j < len(srBuf); j++ {
			sBuf[j] = srBuf[len(srBuf)-j-1]
		}

		cFe.Unpack(&buf)
		cS.Unpack(&sBuf)
		cP.Elligator(&cFe)
		cP2.ScalarMult(&cP, &cS)
		cP2.Pack(&cBuf)

		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
		ep2.ScalarMult(&ep, &sBuf)
		ep2.RistrettoInto(&goBuf)

		if !bytes.Equal(cBuf[:], goBuf[:]) {
			t.Fatalf("%d: %v . %v = %v != %v", i, biS, ep, ep2, cP2)
		}
	}
}

func TestRistrettoEqualsI(t *testing.T) {
	var ep1, ep2 edwards25519.ExtendedPoint
	var torsion [4]edwards25519.ExtendedPoint
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var buf [32]byte
	torsion[0].SetZero()
	torsion[1].SetTorsion1()
	torsion[2].SetTorsion2()
	torsion[3].SetTorsion3()
	for i := 0; i < 1000; i++ {
		rnd.Read(buf[:])
		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep1.SetCompleted(&cp)
		for j := 0; j < 4; j++ {
			ep2.Add(&ep1, &torsion[j])
			if ep1.RistrettoEqualsI(&ep2) != 1 {
				t.Fatalf("%v + %v != %v", ep1, torsion[j], ep2)
			}
		}
	}
}

func TestProjectiveJacobiQuarticConversions(t *testing.T) {
	var buf [32]byte
	var feZero, feOne, fe, js, jt, zInv edwards25519.FieldElement
	var cp, cp2 edwards25519.CompletedPoint
	var ep, ep2 edwards25519.ExtendedPoint
	var jp edwards25519.ProjectiveJacobiPoint
	feOne.SetOne()
	feZero.SetZero()
	for i := 0; i < 1000; i++ {
		if i == 0 {
			ep = edwards25519.ExtendedPoint{feZero, feOne, feOne, feZero}
		} else if i == 1 {
			ep = edwards25519.ExtendedPoint{feOne, feZero, feOne, feZero}
		} else {
			rnd.Read(buf[:])
			fe.SetBytes(&buf)
			cp.SetRistrettoElligator2(&fe)
			ep.SetCompleted(&cp)
		}
		jp.SetExtended(&ep)

		zInv.Inverse(&jp.Z)
		js.Mul(&jp.S, &zInv)
		jt.Mul(&jp.T, &zInv)
		jt.Mul(&jt, &zInv)
		cp2.SetJacobiQuartic(&js, &jt)
		ep2.SetCompleted(&cp2)

		if ep2.RistrettoEqualsI(&ep) != 1 {
			t.Logf("%v", &jp)
			t.Fatalf("Jacobi(Jacobi^-1(%v)) == %v", &ep, &ep2)
		}
	}
}

func TestRistrettoElligator2Inverse(t *testing.T) {
	var buf [32]byte
	var fe edwards25519.FieldElement
	var cp, cp2 edwards25519.CompletedPoint
	var ep, ep2 edwards25519.ExtendedPoint
	var fs [8]edwards25519.FieldElement
	for i := 0; i < 1000; i++ {
		ok := true
		rnd.Read(buf[:])
		buf[31] &= 127
		buf[0] &= 254
		fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
		setMask := ep.RistrettoElligator2Inverse(&fs)
		foundOriginal := false
		count := 0
		for j := 0; j < 8; j++ {
			if ((1 << uint(j)) & setMask) == 0 {
				continue
			}
			if fs[j].Equals(&fe) {
				foundOriginal = true
			}
			cp2.SetRistrettoElligator2(&fs[j])
			ep2.SetCompleted(&cp2)
			if ep2.RistrettoEqualsI(&ep) != 1 {
				t.Logf("%vth preimage %v is wrong: %v", j, &fs[j], &ep2)
				ok = false
			}
			count++
		}
		if !foundOriginal {
			t.Logf("Missing original %v among %d preimage(s)", &fe, count)
			ok = false
		}
		if !ok {
			t.Fatalf("^ see errors above.  fe=%v, ep=%v", hex.EncodeToString(buf[:]), &ep)
		}
	}
}

func BenchmarkElligatorPlusInverse(b *testing.B) {
	var fe edwards25519.FieldElement
    var fs [8]edwards25519.FieldElement
	var ep edwards25519.ExtendedPoint
	var cp edwards25519.CompletedPoint
	var buf [32]byte
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
        rnd.Read(buf[:])
        buf[0] &= 254
        buf[31] &= 127
        fe.SetBytes(&buf)
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
		ep.RistrettoElligator2Inverse(&fs)
	}
}

func BenchmarkElligator(b *testing.B) {
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var buf [32]byte
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cp.SetRistrettoElligator2(&fe)
		ep.SetCompleted(&cp)
	}
}

func BenchmarkRistrettoPack(b *testing.B) {
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var buf [32]byte
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ep.RistrettoInto(&buf)
	}
}

func BenchmarkRistrettoUnpack(b *testing.B) {
	var fe edwards25519.FieldElement
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var buf [32]byte
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	ep.RistrettoInto(&buf)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ep.SetRistretto(&buf)
	}
}

func BenchmarkScalarMult(b *testing.B) {
	var buf, sBuf [32]byte
	var biS big.Int
	var cp edwards25519.CompletedPoint
	var ep edwards25519.ExtendedPoint
	var fe edwards25519.FieldElement
	biS.Rand(rnd, &biL)
	srBuf := biS.Bytes()
	for j := 0; j < len(srBuf); j++ {
		sBuf[j] = srBuf[len(srBuf)-j-1]
	}
	rnd.Read(buf[:])
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ep.ScalarMult(&ep, &sBuf)
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
