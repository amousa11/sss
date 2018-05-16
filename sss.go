package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

// Point on a polynomial in Fp[x] where p is a prime
type Point struct {
	x *big.Int
	y *big.Int
}

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
	secret, points := makeRandomShares(n, m, prime)
	fmt.Println("Secret : ", secret.Text(16))

	fmt.Println("Shares : ")
	for i := 0; i < len(points); i++ {
		fmt.Println(points[i].x.Text(16), "\t", points[i].y.Text(16))
	}

	recoveredSecret := recoverSecret(points[:n], prime)

	fmt.Println("secret recovered from minimum subset of shares", recoveredSecret.Text(16))
}

func makeRandomShares(minimum int, shares int, prime *big.Int) (*big.Int, []*Point) {
	poly := make([]*big.Int, minimum)
	points := make([]*Point, shares)

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
		points[i] = evaluatePoly(poly, point, prime)
	}

	return poly[0], points
}

func recoverSecret(points []*Point, prime *big.Int) *big.Int {
	if len(points) < 2 {
		panic("need at least 2 points")
	}

	return lagrangeInterpolate(big.NewInt(0), points, prime)
}

func lagrangeInterpolate(point *big.Int, points []*Point, prime *big.Int) *big.Int {
	// assuming points are distinct
	sum := big.NewInt(0)
	elems := make(chan *big.Int, len(points))
	go lagrangeInterpolateHelper(elems, points, prime) // concurrently calculate prod x_i / x_i - x_j
	for elem := range elems {
		sum.Add(sum, elem) // sum += f(x_j) * prod x_i / x_i - x_j
	}
	sum.Mod(sum, prime)
	return sum
}

func lagrangeInterpolateHelper(elems chan *big.Int, points []*Point, prime *big.Int) {
	for i := 0; i < len(points); i++ {
		prod := big.NewInt(1)
		for j := 0; j < len(points); j++ {
			if j != i {
				denom := big.NewInt(0).Set(points[j].x) // denom = x_j
				denom.Sub(denom, points[i].x)           // denom = x_j - x_i
				denom.ModInverse(denom, prime)          // denom = 1 / (x_j - x_i)
				x := big.NewInt(0).Set(points[j].x)     // x = x_j
				x.Mul(x, denom)                         // x_j * (x_j - x_i)
				prod.Mul(prod, x)                       // prod *= x_j * (x_j - x_i)
			}
		}
		prod.Mul(prod, points[i].y) // prod *= y
		elems <- prod               // send it back to lagrangeInterpolator
	}
	close(elems) // done calculating, close the channel and exit
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

func evaluatePoly(coeff []*big.Int, x *big.Int, prime *big.Int) *Point {
	point := new(Point)
	y := big.NewInt(0)
	for i := 0; i < len(coeff); i++ {
		eval := big.NewInt(0)
		eval.Set(x)                                 // eval = x
		eval.Exp(eval, big.NewInt(int64(i)), prime) // x ^ i
		eval.Mul(eval, coeff[i])                    // * coeff[i]
		y.Add(y, eval)                              // y += coeff * x ^ i
	}
	y.Mod(y, prime) // y = y mod prime
	point.x = x
	point.y = y
	return point
}
