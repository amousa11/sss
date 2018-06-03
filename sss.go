package sss

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"sss/utils"
	"strconv"
)

func main() {

	fmt.Println(os.Args)
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Errorf("expected more arguments: ./sss minimum shares")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Errorf(err.Error())
	}

	m, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Errorf(err.Error())
	}

	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	fmt.Println("FIELD ORDER =", prime.Text(16))
	fmt.Println("Generating ", n, "shares with threshold", m, " for recovery:")
	secret, points, e := GenerateShares(n, m, prime)
	if e != nil {
		fmt.Errorf(e.Error())
	}
	fmt.Println("Secret : ", secret.Text(16))

	fmt.Println("Shares : ")
	for i := 0; i < len(points); i++ {
		fmt.Println(points[i].X.Text(16), "\t", points[i].Y.Text(16))
	}

	recoveredSecret, e := RecoverSecret(points[:n], prime)

	if e != nil {
		fmt.Errorf(e.Error())
	}

	fmt.Println("secret recovered from minimum subset of shares", recoveredSecret.Text(16))
}

// GenerateShares generates a number of shares which can only be recovered by the minimum number of shares
func GenerateShares(minimum int, shares int, prime *big.Int) (*big.Int, []*utils.Point, error) {
	poly := make([]*big.Int, minimum)
	points := make([]*utils.Point, shares)

	if minimum > shares {
		errors.New("Minimum number of shares specified is greater than the total number of shares")
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
		point := big.NewInt(int64(i + 1))
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
