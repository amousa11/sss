package main

import "fmt"
import "math/big"
import "crypto/rand"

func main() {
	prime := big.NewInt(1)
	prime.Mul(prime, big.NewInt(2)).Exp(prime, big.NewInt(127), big.NewInt(0)).Sub(prime, big.NewInt(1))
	fmt.Println("FIELD ORDER 2^512 - 1 =", prime.Text(16))

	secret, xs, ys := makeRandomShares(3, 6, prime)
	fmt.Println("Secret : ", secret.Text(16))

	fmt.Println("Points : ")
	for i := 0; i < len(xs); i++ {
		fmt.Println(xs[i].Text(16), "\t", ys[i].Text(16))
	}

	recoveredSecret := recoverSecret(xs[:3], ys[:3], prime)

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
	for j := 0; j < len(xs); j++ {
		prod := big.NewInt(1)
		for m := 0; m < len(xs); m++ {
			if m != j {
				denom := big.NewInt(0).Set(xs[m])
				denom.Sub(denom, xs[j])
				denom.ModInverse(denom, prime)
				x := big.NewInt(0).Set(xs[m])
				x.Mul(x, denom)
				prod.Mul(prod, x)
			}
		}
		y := big.NewInt(0).Set(ys[j])
		y.Mul(y, prod)
		sum.Add(sum, y)
	}
	sum.Mod(sum, prime)
	return sum
}
