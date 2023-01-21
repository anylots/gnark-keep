package main

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// Circuit defines a simple circuit
// x**3 + x + 5 == y
type CircuitScoreUp struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x > y
func (circuit *CircuitScoreUp) Define(api frontend.API) error {
	// Cmp returns 1 if i1>i2, 0 if i1=i2, -1 if i1<i2
	result := api.Cmp(circuit.X, circuit.Y)
	api.AssertIsEqual(result, 1)
	return nil
}

func main() {
	//step1. instantiate circuit
	var circuitScoreUp CircuitScoreUp
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuitScoreUp)
	if err != nil {
		return
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return
	}

	//step2. export the groth16.VerifyingKey as a solidity smart contract.
	var fileName = "./verify/verifyScoreUp.sol"
	var solidityFile, _ = os.Create(fileName)
	// writer := bufio.NewWriter(solidityFile)
	vk.ExportSolidity(solidityFile)

	//step3. generate witness and prove
	assignment := &CircuitScoreUp{
		X: frontend.Variable(2),
		Y: frontend.Variable(1),
	}

	witness, _ := frontend.NewWitness(assignment, ecc.BN254)

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return
	}

	//step4. generate public witness and verify
	validPublicWitness, _ := frontend.NewWitness(assignment, ecc.BN254, frontend.PublicOnly())
	err = groth16.Verify(proof, vk, validPublicWitness)
	if err != nil {
		fmt.Printf("verification failed\n")
		return
	}
	fmt.Printf("verification succeded\n")
}

function prove(assignment: &CircuitCubic){
	witness, _ := frontend.NewWitness(assignment, ecc.BN254)

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return
	}
}
