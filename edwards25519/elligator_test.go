package edwards25519_test

import (
	"bytes"
	"encoding/hex"
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

func BenchmarkElligatorInverse(b *testing.B) {
	var fe edwards25519.FieldElement
	var fs [8]edwards25519.FieldElement
	var ep edwards25519.ExtendedPoint
	var cp edwards25519.CompletedPoint
	var buf [32]byte
	rnd.Read(buf[:])
	buf[0] &= 254
	buf[31] &= 127
	fe.SetBytes(&buf)
	cp.SetRistrettoElligator2(&fe)
	ep.SetCompleted(&cp)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
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
