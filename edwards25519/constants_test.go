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
}
