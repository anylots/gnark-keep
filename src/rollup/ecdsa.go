package main

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)


// Circuit "toy" rollup circuit where an operator can generate a proof that he processed
// some transactions
type AccountCircuit struct {
	// ---------------------------------------------------------------------------------------------
	// SECRET INPUTS

	// list of accounts involved before update and their public keys
	SenderAccountsBefore   [batchSize]AccountConstraints
	ReceiverAccountsBefore [batchSize]AccountConstraints

	// list of accounts involved after update and their public keys
	SenderAccountsAfter   [batchSize]AccountConstraints
	ReceiverAccountsAfter [batchSize]AccountConstraints

	// list of transactions
	Transfers [batchSize]TransferConstraints
}