package edwards25519

import (
	"testing"
)

func TestFeConstants(t *testing.T) {
	var dP1, dP1InvDM1 FieldElement
	dP1.add(&feD, &feOne)
	dP1InvDM1.sub(&feD, &feOne)
	dP1InvDM1.Inverse(&dP1InvDM1)
	dP1InvDM1.Mul(&dP1InvDM1, &dP1)
	if dP1InvDM1.EqualsI(&feDp1OverDm1) != 1 {
		t.Fatalf("feDp1OverDm1: %v != %v", &feDp1OverDm1, &dP1InvDM1)
	}

	var sqrtID FieldElement
	sqrtID.Mul(&feI, &feD)
	sqrtID.Sqrt(&sqrtID)
	if sqrtID.EqualsI(&feSqrtID) != 1 {
		t.Fatalf("sqrtID: %v != %v", &feSqrtID, &sqrtID)
	}

	var doubleInvSqrtMinusDMinusOne FieldElement
	doubleInvSqrtMinusDMinusOne.Add(&feInvSqrtMinusDMinusOne, &feInvSqrtMinusDMinusOne)
	if doubleInvSqrtMinusDMinusOne.EqualsI(&feDoubleInvSqrtMinusDMinusOne) != 1 {
		t.Fatalf("doubleInvSqrtMinusDMinusOne: %v != %v",
			feDoubleInvSqrtMinusDMinusOne, doubleInvSqrtMinusDMinusOne)
	}

	var doubleIInvSqrtMinusDMinusOne FieldElement
	doubleIInvSqrtMinusDMinusOne.Mul(&feDoubleInvSqrtMinusDMinusOne, &feI)
	if doubleIInvSqrtMinusDMinusOne.EqualsI(&feDoubleIInvSqrtMinusDMinusOne) != 1 {
		t.Fatalf("doubleIInvSqrtMinusDMinusOne: %v != %v",
			feDoubleIInvSqrtMinusDMinusOne, doubleIInvSqrtMinusDMinusOne)
	}

	var invSqrt1pD FieldElement
	invSqrt1pD.add(&feD, &feOne)
	invSqrt1pD.InvSqrt(&invSqrt1pD)
	if invSqrt1pD.EqualsI(&feInvSqrt1pD) != 1 {
		t.Fatalf("invSqrt1pD: %v != %v", feInvSqrt1pD, invSqrt1pD)
	}
}
