package main

import (
	"fmt"
	"os"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

)

const (
	nbAccounts = 16 // 16 accounts so we know that the proof length is 5
	depth      = 5  // size fo the inclusion proofs
	batchSize  = 1  // nbTranfers to batch in a proof
)

// Circuit "toy" rollup circuit where an operator can generate a proof that he processed
// some transactions
type RollupCircuit struct {
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

// AccountConstraints accounts encoded as constraints
type AccountConstraints struct {
	Index   frontend.Variable // index in the tree
	Nonce   frontend.Variable // nb transactions done so far from this account
	Balance frontend.Variable `gnark:",public"`
}

// TransferConstraints transfer encoded as constraints
type TransferConstraints struct {
	Amount         frontend.Variable `gnark:",public"`
}

// func (circuit *RollupCircuit) postInit(api frontend.API) error {

// 	for i := 0; i < batchSize; i++ {

// 		// setting the transfers
// 		circuit.Transfers[i].Nonce = circuit.SenderAccountsBefore[i].Nonce
// 		circuit.Transfers[i].SenderPubKey = circuit.PublicKeysSender[i]
// 		circuit.Transfers[i].ReceiverPubKey = circuit.PublicKeysReceiver[i]

// 	}
// 	return nil
// }

// Define declares the circuit's constraints
func (circuit *RollupCircuit) Define(api frontend.API) error {
	// if err := circuit.postInit(api); err != nil {
	// 	return err
	// }

	// creation of the circuit
	for i := 0; i < batchSize; i++ {

		//ecdsa

		//merkle

		// update the accounts
		verifyAccountUpdated(api, circuit.SenderAccountsBefore[i], circuit.ReceiverAccountsBefore[i], 
			circuit.SenderAccountsAfter[i], circuit.ReceiverAccountsAfter[i], circuit.Transfers[i].Amount)
	}

	return nil
}

func verifyAccountUpdated(api frontend.API, from, to, fromUpdated, toUpdated AccountConstraints, amount frontend.Variable) {

	// ensure that nonce is correctly updated
	nonceUpdated := api.Add(from.Nonce, 1)
	api.AssertIsEqual(nonceUpdated, fromUpdated.Nonce)

	nonceUpdated_to := api.Add(to.Nonce, 1)
	api.AssertIsEqual(nonceUpdated_to, toUpdated.Nonce)


	// ensure that account is correctly
	api.AssertIsEqual(from.Index, fromUpdated.Index)
	api.AssertIsEqual(to.Index, toUpdated.Index)


	// ensures that the amount is less than the balance
	api.AssertIsLessOrEqual(amount, from.Balance)

	// ensure that balance is correctly updated
	fromBalanceUpdated := api.Sub(from.Balance, amount)
	api.AssertIsEqual(fromBalanceUpdated, fromUpdated.Balance)

	toBalanceUpdated := api.Add(to.Balance, amount)
	api.AssertIsEqual(toBalanceUpdated, toUpdated.Balance)

}

func main() {
	//step1. instantiate circuit
	var rollupCircuit RollupCircuit
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &rollupCircuit, frontend.IgnoreUnconstrainedInputs())
	if err != nil {
		return
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return
	}

	// step2. export the groth16.VerifyingKey as a solidity smart contract.
	var fileName = "./verify/verifyRollup.sol"
	var solidityFile, _ = os.Create(fileName)
	// writer := bufio.NewWriter(solidityFile)
	vk.ExportSolidity(solidityFile)

	//step3. generate witness and prove
	var assignment RollupCircuit

	// set witnesses for the accounts before update
	assignment.SenderAccountsBefore[0].Index = 0
	assignment.SenderAccountsBefore[0].Nonce = 0
	assignment.SenderAccountsBefore[0].Balance = 100

	assignment.ReceiverAccountsBefore[0].Index = 1
	assignment.ReceiverAccountsBefore[0].Nonce = 0
	assignment.ReceiverAccountsBefore[0].Balance = 0

	//// set witnesses for the transfer
	assignment.Transfers[0].Amount = 20

    // set the witnesses for the account after update
	assignment.SenderAccountsAfter[0].Index = 0
	assignment.SenderAccountsAfter[0].Nonce = 1
	assignment.SenderAccountsAfter[0].Balance = 80

	assignment.ReceiverAccountsAfter[0].Index = 1
	assignment.ReceiverAccountsAfter[0].Nonce = 1
	assignment.ReceiverAccountsAfter[0].Balance = 20

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254)

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return
	}

	//step4. generate public witness and verify

	// assignment.ReceiverAccountsAfter[0].Balance = 30
	validPublicWitness, _ := frontend.NewWitness(&assignment, ecc.BN254, frontend.PublicOnly())
	err = groth16.Verify(proof, vk, validPublicWitness)
	if err != nil {
		fmt.Printf("verification failed\n")
		return
	}
	fmt.Printf("verification succeded\n")
}
