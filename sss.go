package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func main() {

	fmt.Println(os.Args)
	args := os.Args[1:]

	if len(args) < 2 {
		panic("expected more arguments: ./sss minimum shares")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}

	m, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	fmt.Println("FIELD ORDER =", prime.Text(16))
	fmt.Println("Generating ", n, "shares with threshold", m, " for recovery:")
	secret, xs, ys := makeRandomShares(n, m, prime)
	fmt.Println("Secret : ", secret.Text(16))

	fmt.Println("Shares : ")
	for i := 0; i < len(xs); i++ {
		fmt.Println(xs[i].Text(16), "\t", ys[i].Text(16))
	}

	recoveredSecret := recoverSecret(xs[:n], ys[:n], prime)

	// fmt.Println("secret", secret.Text(16))

	fmt.Println("secret recovered from minimum subset of shares", recoveredSecret.Text(16))
}

func makeRandomShares(minimum int, shares int, prime *big.Int) (*big.Int, []*big.Int, []*big.Int) {
	xs := make([]*big.Int, shares)
	ys := make([]*big.Int, shares)
	poly := make([]*big.Int, minimum)

	if minimum > shares {
		panic("min less than shares")
	}

	for i := 0; i < minimum; i++ { // should be i < shares.
		coeff := generateRandomBytes(32)
		coeff.Mod(coeff, prime)
		poly[i] = coeff
	}

	for i := 0; i < shares; i++ {
		point := big.NewInt(int64(i + 1))
		xs[i], ys[i] = evaluatePoly(poly, point, prime)
	}

	return poly[0], xs, ys
}

func generateRandomBytes(n int) *big.Int {
	b := make([]byte, n)
	n, err := rand.Read(b)

	if err != nil {
		panic("rand byte generation error")
	}

	rtn := big.NewInt(0)
	rtn.SetBytes(b)
	return rtn
}

func recoverSecret(xs []*big.Int, ys []*big.Int, prime *big.Int) *big.Int {
	if len(xs) < 2 {
		panic("need at least 2 points")
	}

	return lagrangeInterpolate(big.NewInt(0), xs, ys, prime)
}

func evaluatePoly(coeff []*big.Int, point *big.Int, prime *big.Int) (*big.Int, *big.Int) {
	total := big.NewInt(0)
	for i := 0; i < len(coeff); i++ {
		x := big.NewInt(0)
		x.Set(point)
		x.Exp(x, big.NewInt(int64(i)), prime)
		x.Mul(x, coeff[i])
		total.Add(total, x)
	}
	total.Mod(total, prime)
	return point, total
}

func lagrangeInterpolate(point *big.Int, xs []*big.Int, ys []*big.Int, prime *big.Int) *big.Int {
	// assuming points are distinct
	sum := big.NewInt(0)
	prod := big.NewInt(1)
	elems := make(chan *big.Int, len(xs))
	for j := 0; j < len(xs); j++ {
		go lagrangeInterpolateHelper(elems, xs, j, prime)
		y := big.NewInt(0).Set(ys[j])
		y.Mul(y, prod)
		sum.Add(sum, y)
	}
	sum.Mod(sum, prime)
	return sum
}

func lagrangeInterpolateHelper(elems chan *big.Int, arr []*big.Int, ind int, prime *big.Int) {
	for m := 0; m < len(arr); m++ {
		if m != ind {
			denom := big.NewInt(0).Set(arr[m])
			denom.Sub(denom, arr[ind])
			denom.ModInverse(denom, prime)
			x := big.NewInt(0).Set(arr[m])
			x.Mul(x, denom)
			elems <- x
		}
	}
}
