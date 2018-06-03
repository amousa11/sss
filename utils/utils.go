package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Point is a struct containing X and Y, two *big.Int which represent a point in a field F_m
type Point struct {
	X *big.Int
	Y *big.Int
}

func (p Point) String() string {
	return p.X.String() + "\t" + p.Y.String()
}

// LagrangeInterpolate performs lagrange interpolation on an array of *Points and returns a value
func LagrangeInterpolate(point *big.Int, points []*Point, modulus *big.Int) *big.Int {
	// assuming points are distinct
	sum := big.NewInt(0)
	elems := make(chan *big.Int, len(points))
	go lagrangeInterpolateHelper(elems, points, modulus) // concurrently calculate prod x_i / x_i - x_j
	for elem := range elems {
		sum.Add(sum, elem) // sum += f(x_j) * prod x_i / x_i - x_j
	}
	sum.Mod(sum, modulus)
	return sum
}

// Subroutine of LagrangeInterpolate
func lagrangeInterpolateHelper(elems chan *big.Int, points []*Point, modulus *big.Int) {
	for i := 0; i < len(points); i++ {
		prod := big.NewInt(1)
		for j := 0; j < len(points); j++ {
			if j != i {
				denom := big.NewInt(0).Set(points[j].X) // denom = x_j
				denom.Sub(denom, points[i].X)           // denom = x_j - x_i
				denom.ModInverse(denom, modulus)        // denom = 1 / (x_j - x_i)
				x := big.NewInt(0).Set(points[j].X)     // x = x_j
				x.Mul(x, denom)                         // x_j * (x_j - x_i)
				prod.Mul(prod, x)                       // prod *= x_j * (x_j - x_i)
			}
		}
		prod.Mul(prod, points[i].Y) // prod *= y
		elems <- prod               // send it back to lagrangeInterpolator
	}
	close(elems) // done calculating, close the channel and exit
}

// GenerateRandomBigInt generates n random bytes and returns it as a *big.Int
func GenerateRandomBigInt(n int) (*big.Int, error) {
	b := make([]byte, n)
	n, err := rand.Read(b)

	if err != nil {
		return nil, errors.New("rand byte generation error")
	}

	rtn := big.NewInt(0)
	rtn.SetBytes(b)
	return rtn, nil
}

// EvaluatePolynomial evaluates a polynomial at a certain xValue under a modulus and returns a *Point
func EvaluatePolynomial(coefficients []*big.Int, xValue *big.Int, modulus *big.Int) *Point {
	point := new(Point)
	yValue := big.NewInt(0)
	for i := 0; i < len(coefficients); i++ {
		eval := big.NewInt(0)
		eval.Set(xValue)                              // eval = x
		eval.Exp(eval, big.NewInt(int64(i)), modulus) // x ^ i
		eval.Mul(eval, coefficients[i])               // * coeff[i]
		yValue.Add(yValue, eval)                      // y += coeff * x ^ i
	}
	yValue.Mod(yValue, modulus) // y = y mod prime
	point.X = xValue
	point.Y = yValue
	return point
}
