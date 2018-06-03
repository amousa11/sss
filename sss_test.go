package sss

import (
	"math/big"
	"testing"
)

func TestShareGenerationAndRecovery(t *testing.T) {

	min := 200
	shares := 500
	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	secret, points, e1 := GenerateShares(min, shares, prime)
	recoveredSecret, e2 := RecoverSecret(points[:min], prime)

	if e1 != nil {
		t.Error(e1)
		t.Fail()
	}

	if e2 != nil {
		t.Error(e2)
		t.Fail()
	}

	t.Log("Expect ", secret, "to equal", recoveredSecret)

	if secret.Cmp(recoveredSecret) != 0 {
		t.Errorf("Expected secret %x to equal recovered secret %x", secret, recoveredSecret)
	}

}

func benchmarkShareGeneration(i int, b *testing.B) {
	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GenerateShares(i/3, i, prime)
	}
}

func BenchmarkShareGeneration10(b *testing.B) {
	benchmarkShareGeneration(10, b)
}

func BenchmarkShareGeneration20(b *testing.B) {
	benchmarkShareGeneration(20, b)
}

func BenchmarkShareGeneration50(b *testing.B) {
	benchmarkShareGeneration(50, b)
}

func BenchmarkShareGeneration100(b *testing.B) {
	benchmarkShareGeneration(100, b)
}

func BenchmarkShareGeneration200(b *testing.B) {
	benchmarkShareGeneration(200, b)
}

func BenchmarkShareGeneration500(b *testing.B) {
	benchmarkShareGeneration(500, b)
}

func BenchmarkShareGeneration1000(b *testing.B) {
	benchmarkShareGeneration(1000, b)
}

func benchmarkSecretRecovery(i int, b *testing.B) {
	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	_, shares, _ := GenerateShares(i/3, i, prime)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		RecoverSecret(shares, prime)
	}
}

func BenchmarkSecretRecovery10(b *testing.B) {
	benchmarkSecretRecovery(10, b)
}
func BenchmarkSecretRecovery20(b *testing.B) {
	benchmarkSecretRecovery(20, b)
}

func BenchmarkSecretRecovery50(b *testing.B) {
	benchmarkSecretRecovery(50, b)
}
func BenchmarkSecretRecovery100(b *testing.B) {
	benchmarkSecretRecovery(100, b)
}
func BenchmarkSecretRecovery200(b *testing.B) {
	benchmarkSecretRecovery(100, b)
}

func BenchmarkSecretRecovery500(b *testing.B) {
	benchmarkSecretRecovery(500, b)
}

func BenchmarkSecretRecovery1000(b *testing.B) {
	benchmarkSecretRecovery(1000, b)
}
