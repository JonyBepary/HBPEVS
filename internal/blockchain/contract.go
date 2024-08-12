package blockchain

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

/*
func (Contract *ContractStruct) GenerateContract(db *leveldb.DB) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(&Contract); err != nil {
		panic(err)
	}

	// generate digest
	var bit []byte
	Contract.Digest = make([]*Digest, len(Contract.Command))
	for i := 0; i < len(Contract.Command); i++ {
		Contract.Digest[i].Digest = Contract.Command[i].GenerateDigest()
	}

	// generate proof of command
	Contract.ProofOfCommand = make([]*ProofOfCommand, len(Contract.Command))
	for i := 0; i < len(Contract.Command); i++ {
		bit = utils.GetAuthorize(Contract.Digest[i].Digest, Contract.Command[i].AuthorizeBy, db)
		Contract.ProofOfCommand[i].ProofOfCommand = utils.BlobhashGeneration(bit)
		err = db.Put([]byte(Contract.ProofOfCommand[i].ProofOfCommand), bit, nil)
		if err != nil {
			panic(err)
		}
	}

}
*/

func (Contract *ContractStruct) assign_PO(params ...string) {

}
func (Contract *ContractStruct) assign_RO(params ...string) {

}
func (Contract *ContractStruct) assign_V(params ...string) {

}
func (Contract *ContractStruct) assign_C(params ...string) {

}
func (Contract *ContractStruct) assign_P(params ...string) {

}
func (Contract *ContractStruct) assign_D(params ...string) {

}
func (Contract *ContractStruct) assign_S(params ...string) {

}
func (Contract *ContractStruct) assign_T(params ...string) {

}
func (Contract *ContractStruct) assign_H(params ...string) {

}
func (Contract *ContractStruct) assign_K(params ...string) {

}
func (Contract *ContractStruct) assign_POC(params ...string) {

}
func (Contract *ContractStruct) assign_Digest(params ...string) {

}

// takes primary keys of returning officer AKS, and db
// check all the command and it's proof through signature verification
// This function accepts the primary keys of the returning officer AKS (Asymmetric  Key Storage) and the corresponding database.
// then performs a thorough check on all the commands and their respective proofs using signature verification.

// This ensures that the integrity and authenticity of the commands are maintained and
// that only authorized personnel are able to access and manipulate the data in the database.
func (contract *ContractStruct) VerifyContract(db *leveldb.DB) error {
	// Verify AuthorizeBy field
	for _, cmd := range contract.Command {
		if !isValidAddress(cmd.AuthorizeBy, db) {
			return fmt.Errorf("invalid AuthorizeBy address: %s", cmd.AuthorizeBy)
		}
	}
	// Verify Entity field
	for _, cmd := range contract.Command {
		if !isValidEntity(cmd.Entity) {
			return fmt.Errorf("invalid Entity: %s", cmd.Entity)
		}
	}

	// Verify Action field
	for _, cmd := range contract.Command {
		if !isValidAction(cmd.Action, db) {
			return fmt.Errorf("invalid Action: %s", cmd.Action)
		}
	}

	// Verify Digest field
	for i, digest := range contract.Digest {
		expectedDigest := contract.Command[i].GenerateDigest()
		if digest.Digest != expectedDigest {
			return fmt.Errorf("invalid Digest: %s, expected: %s", digest.Digest, expectedDigest)
		}
	}

	// Verify ProofOfCommand field
	for i, proof := range contract.ProofOfCommand {
		authorizeBy := contract.Command[i].AuthorizeBy
		digest := contract.Digest[i].Digest
		bit := utils.GetAuthorize(digest, authorizeBy, db)
		expectedProof := utils.BlobhashGeneration(bit)

		if proof.ProofOfCommand != expectedProof {
			return fmt.Errorf("invalid ProofOfCommand: %s, expected: %s", proof.ProofOfCommand, expectedProof)
		}
	}

	return nil
}

func (vContract *ContractStruct) getDigest() string {
	digest := sha256.New()
	for i := 0; i < len(vContract.Digest); i++ {
		digest.Write([]byte(vContract.Digest[i].Digest))
	}
	return base58.Encode(digest.Sum(nil))
}

func (c ContractStruct) show() {
	for i := 0; i < len(c.Command); i++ {
		fmt.Print("Contract ID: ", i, "\n")
		fmt.Print("AuthorizeBy: ", c.Command[i].AuthorizeBy, "\n")
		fmt.Print("Command: ", c.Command[i].Action, "\n")
		fmt.Print("Entity: ", c.Command[i].Entity, "\n")
		fmt.Print("Digest: ", c.Digest[i].Digest, "\n")
		fmt.Print("Proof: ", c.ProofOfCommand[i].ProofOfCommand, "\n")
	}

}
func isValidAddress(address string, db *leveldb.DB) bool {
	// Check if the address exists on the blockchain
	dbAddress, err := db.Get([]byte("isValidAddress"), nil)
	if err != nil {
		return false
	}

	var Set *treeset.Set
	err = Set.UnmarshalJSON(dbAddress)
	if err != nil {
		log.Panic(err)
		return false
	}
	return Set.Contains(address)
}

func isValidAction(action string, db *leveldb.DB) bool {
	// Define a list of valid actions
	validActions, err := db.Get([]byte("isValidAction"), nil)

	if err != nil {
		log.Panic(err)
	}
	var Set *treeset.Set
	err = Set.UnmarshalJSON(validActions)
	if err != nil {
		log.Panic(err)
	}
	if Set.Contains(action) {
		return true
	}

	return false
}

func isValidEntity(entity string) bool {
	// Check if format of the entity is valid
	ln := len(entity)
	if ln == 64 || ln == 44 {
		return true
	}
	return false
}

func suspend_PO(entity, authorizeBy string, db *leveldb.DB) bool {

	contractFile, err := db.Get([]byte("contract"), nil)
	if err != nil {
		log.Panic(err)
		return false
	}
	var c *ContractStruct
	c.XXX_Unmarshal(contractFile)

	// append new command
	c.Command = append(c.Command, &Command{AuthorizeBy: authorizeBy, Action: "SUSPEND", Entity: entity})
	c.Digest = append(c.Digest, &Digest{Digest: utils.String_hash_generation(c.Command[len(c.Command)-1].GenerateDigest())})
	c.ProofOfCommand = append(c.ProofOfCommand, &ProofOfCommand{
		ProofOfCommand: utils.BlobhashGeneration(utils.GetAuthorize(c.Digest[len(c.Digest)-1].Digest, authorizeBy, db))})

	c.show()
	//PUT THE CONTRACT BACK TO THE DB
	contractFile, err = c.XXX_Marshal(nil, false)
	if err != nil {
		log.Panic(err)
		return false
	}
	err = db.Put([]byte("contract"), contractFile, nil)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}
func assign_PO(entity, authorizeBy string, db *leveldb.DB) bool {

	contractFile, err := db.Get([]byte("contract"), nil)
	if err != nil {
		log.Panic(err)
		return false
	}
	var c *ContractStruct
	c.XXX_Unmarshal(contractFile)

	// append new command
	c.Command = append(c.Command, &Command{AuthorizeBy: authorizeBy, Action: "ASSIGN", Entity: entity})
	c.Digest = append(c.Digest, &Digest{Digest: utils.String_hash_generation(c.Command[len(c.Command)-1].GenerateDigest())})
	c.ProofOfCommand = append(c.ProofOfCommand, &ProofOfCommand{
		ProofOfCommand: utils.BlobhashGeneration(utils.GetAuthorize(c.Digest[len(c.Digest)-1].Digest, authorizeBy, db))})

	c.show()
	//PUT THE CONTRACT BACK TO THE DB
	contractFile, err = c.XXX_Marshal(nil, false)
	if err != nil {
		log.Panic(err)
		return false
	}
	err = db.Put([]byte("contract"), contractFile, nil)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

func (c *Command) Set(AuthorizeBy, Action, Entity string) {
	//AuthorizeBy Action  Entity
	c.AuthorizeBy = AuthorizeBy
	c.Action = Action
	c.Entity = Entity

}

func (c *Command) GenerateDigest() string {
	digest := sha256.New()
	digest.Write([]byte(c.AuthorizeBy))
	digest.Write([]byte(c.Action))
	digest.Write([]byte(c.Entity))

	//RERURN DIGEST
	return base58.Encode(digest.Sum(nil))

}
