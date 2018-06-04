package sss

import (
	"errors"
	"math/big"

	"github.com/amousa11/sss/utils"
)

// GenerateShares generates a number of shares which can only be recovered by the minimum number of shares
func GenerateShares(minimum int, shares int, prime *big.Int) (*big.Int, []*utils.Point, error) {
	poly := make([]*big.Int, minimum)
	points := make([]*utils.Point, shares)

	if minimum > shares {
		return nil, nil, errors.New("Minimum number of shares specified is greater than the total number of shares")
	}

	if minimum < 2 {
		return nil, nil, errors.New("Minimum number of shares specified is greater than the total number of shares")
	}

	for i := 0; i < minimum; i++ { // should be i < shares.
		coefficients, e := utils.GenerateRandomBigInt(32)
		if e != nil {
			return nil, nil, e
		}
		coefficients.Mod(coefficients, prime)
		poly[i] = coefficients
	}

	for i := 0; i < shares; i++ {
		randXValue, e := utils.GenerateRandomBigInt(32)
		if e != nil {
			return nil, nil, e
		}
		point := randXValue
		points[i] = utils.EvaluatePolynomial(poly, point, prime)
	}

	return poly[0], points, nil
}

// RecoverSecret recovers a secret given an array of *utils.Points and a prime modulus for the Field the points reside in
func RecoverSecret(points []*utils.Point, prime *big.Int) (*big.Int, error) {
	if len(points) < 2 {
		return nil, errors.New("Requires at least 2 shares to recover a secret")
	}

	return utils.LagrangeInterpolate(big.NewInt(0), points, prime), nil
}
