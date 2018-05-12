package main

import "testing"
import "math/big"

func TestSecretGeneration(t *testing.T) {

	min := 200
	shares := 500
	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	secret, xs, ys := makeRandomShares(min, shares, prime)
	recoveredSecret := recoverSecret(xs[:min], ys[:min], prime)

	t.Log("Expect ", secret, "to equal", recoveredSecret)

	if secret.Cmp(recoveredSecret) != 0 {
		t.Errorf("Expected secret %x to equal recovered secret %x", secret, recoveredSecret)
	}

}
