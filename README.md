# sss
Shamir Secret Sharing 

A simple, lightly tested library for share generation and recovery of a secret. 

## Documentation

### GenerateShares

`GenerateShares(minimum int, shares int, prime *big.Int) (*big.Int, []*utils.Point, error)`

This function creates a set of shares returned as an array of points, as well as the secret that these shares recover as a big.Int. It also returns an error.

### RecoverSecret

`RecoverSecret(points []*utils.Point, modulus *big.Int) (*big.Int, error)`

This function recovers a secret from a set of points under `prime` modulus. The secret is returned as a big.Int