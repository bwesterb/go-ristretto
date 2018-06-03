package ristretto_test

import (
	"math/big"
	"math/rand"
	"os"
	"testing"

	"github.com/bwesterb/go-ristretto"
)

var biL big.Int
var rnd *rand.Rand

func TestScPacking(t *testing.T) {
	var bi big.Int
	var s, s2 ristretto.Scalar
	var buf [32]byte
	for i := 0; i < 100; i++ {
		bi.Rand(rnd, &biL)
		s.SetBigInt(&bi)
		s.BytesInto(&buf)
		s2.SetBytes(&buf)
		if s.BigInt().Cmp(s2.BigInt()) != 0 {
			t.Fatalf("Unpack o Pack != id (%v != %v)", &bi, s2.BigInt())
		}
	}
}

func TestScBigIntPacking(t *testing.T) {
	var bi big.Int
	var s ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi.Rand(rnd, &biL)
		s.SetBigInt(&bi)
		if s.BigInt().Cmp(&bi) != 0 {
			t.Fatalf("BigInt o SetBigInt != id (%v != %v)", &bi, s.BigInt())
		}
	}
}

func TestScSub(t *testing.T) {
	var bi1, bi2, bi3 big.Int
	var s1, s2, s3 ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi1.Rand(rnd, &biL)
		bi2.Rand(rnd, &biL)
		bi3.Sub(&bi1, &bi2)
		bi3.Mod(&bi3, &biL)
		s1.SetBigInt(&bi1)
		s2.SetBigInt(&bi2)
		if s3.Sub(&s1, &s2).BigInt().Cmp(&bi3) != 0 {
			t.Fatalf("%v - %v = %v != %v", &bi1, &bi2, &bi3, s3.BigInt())
		}
	}
}

func TestScAdd(t *testing.T) {
	var bi1, bi2, bi3 big.Int
	var s1, s2, s3 ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi1.Rand(rnd, &biL)
		bi2.Rand(rnd, &biL)
		bi3.Add(&bi1, &bi2)
		bi3.Mod(&bi3, &biL)
		s1.SetBigInt(&bi1)
		s2.SetBigInt(&bi2)
		if s3.Add(&s1, &s2).BigInt().Cmp(&bi3) != 0 {
			t.Fatalf("%v + %v = %v != %v", &bi1, &bi2, &bi3, s3.BigInt())
		}
	}
}

func TestScMul(t *testing.T) {
	var bi1, bi2, bi3 big.Int
	var s1, s2, s3 ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi1.Rand(rnd, &biL)
		bi2.Rand(rnd, &biL)
		bi3.Mul(&bi1, &bi2)
		bi3.Mod(&bi3, &biL)
		s1.SetBigInt(&bi1)
		s2.SetBigInt(&bi2)
		if s3.Mul(&s1, &s2).BigInt().Cmp(&bi3) != 0 {
			t.Fatalf("%v * %v = %v != %v", &bi1, &bi2, &bi3, s3.BigInt())
		}
	}
}

func TestScMulAdd(t *testing.T) {
	var bi1, bi2, bi3, bi4 big.Int
	var s1, s2, s3, s4 ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi1.Rand(rnd, &biL)
		bi2.Rand(rnd, &biL)
		bi3.Rand(rnd, &biL)
		bi4.Mul(&bi1, &bi2)
		bi4.Add(&bi4, &bi3)
		bi4.Mod(&bi4, &biL)
		s1.SetBigInt(&bi1)
		s2.SetBigInt(&bi2)
		s3.SetBigInt(&bi3)
		if s4.MulAdd(&s1, &s2, &s3).BigInt().Cmp(&bi4) != 0 {
			t.Fatalf("%v * %v + %v = %v != %v",
				&bi1, &bi2, &bi3, &bi4, s4.BigInt())
		}
	}
}

func TestScInverse(t *testing.T) {
	var bi1, bi2 big.Int
	var s1, s2 ristretto.Scalar
	for i := 0; i < 100; i++ {
		bi1.Rand(rnd, &biL)
		bi2.ModInverse(&bi1, &biL)
		s1.SetBigInt(&bi1)
		if s2.Inverse(&s1).BigInt().Cmp(&bi2) != 0 {
			t.Fatalf("1/%v = %v != %v", &bi1, &bi2, &s2)
		}
	}
}

func TestMain(m *testing.M) {
	biL.SetString(
		"1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed", 16)
	rnd = rand.New(rand.NewSource(37))
	os.Exit(m.Run())
}
