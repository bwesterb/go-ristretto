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
		t.Fatalf("feDp1OverDm1: %v != %v", &feDp1OverDm1, dP1InvDM1)
	}
}
