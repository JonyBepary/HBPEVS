package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cloudflare/circl/sign/schemes"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

// will interfaced to ReturningOfficer struct
func (v *ReturningOfficer) GenerateDigest() []byte {

	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", v.Name)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Id)))
	for _, vr := range v.GetRArea() {
		digest.Write([]byte(fmt.Sprintf("%v", vr)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.PublicKey)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.EncryptionAlgorithm)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.KeySize)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.SigSize)))

	return digest.Sum(nil)
}

func (v *ReturningOfficer) Verify(primaryKey string, db *leveldb.DB) bool {
	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", v.Name)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Id)))
	for _, vr := range v.GetRArea() {
		digest.Write([]byte(fmt.Sprintf("%v", vr)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.PublicKey)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.EncryptionAlgorithm)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.KeySize)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.SigSize)))
	blob := utils.GetAuthorize(base58.Encode(digest.Sum(nil)), primaryKey, db)
	return v.GetSignature() == utils.BlobhashGeneration(blob)
}

// Generate Digest of ReturningOfficer struct
func (v *PollingOfficer) GenerateDigest() []byte {
	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", v.Name)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Id)))
	digest.Write([]byte(fmt.Sprintf("%v", v.PollingStation)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.PublicKey)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.EncryptionAlgorithm)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.KeySize)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.SigSize)))
	return digest.Sum(nil)
}

func (v *PollingOfficer) Verify(primaryKey string, db *leveldb.DB) bool {
	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", v.Name)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Id)))
	digest.Write([]byte(fmt.Sprintf("%v", v.PollingStation)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.PublicKey)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.EncryptionAlgorithm)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.KeySize)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.SigSize)))
	blob := utils.GetAuthorize(base58.Encode(digest.Sum(nil)), primaryKey, db)
	return v.GetSignature() == utils.BlobhashGeneration(blob)
}

// will interfaced to ReturningOfficer struct
func GenerateReturningOfficer(AuthorizeBy string, db *leveldb.DB) []ReturningOfficer {
	var ReturningOfficers []ReturningOfficer

	//parse info from configuration
	configFile, err := os.Open("config/returning_officer.json")
	check(err)
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&ReturningOfficers); err != nil {
		panic(err)
	}
	// spew.Dump(ReturningOfficers)

	var blob []byte
	for i := 0; i < len(ReturningOfficers); i++ {
		//generating aks
		ReturningOfficers[i].Keys = &AKS{}
		ReturningOfficers[i].Keys.generate(db, path.Join("internal/role/RO/", ReturningOfficers[i].Id))
		hash := ReturningOfficers[i].GenerateDigest()
		ReturningOfficers[i].Digest = base58.Encode(hash)
		blob = utils.GetAuthorize(ReturningOfficers[i].Digest, AuthorizeBy, db)
		ReturningOfficers[i].Signature = utils.BlobhashGeneration(blob)
		db.Put([]byte(ReturningOfficers[i].Signature), blob, nil)

	}
	return ReturningOfficers
}

// will interfaced to ReturningOfficer struct
func GeneratePollingOfficer(AuthorizeBy string, db *leveldb.DB) []PollingOfficer {

	var PO []PollingOfficer

	//parse info from configuration
	configFile, err := os.Open("config/polling_officer.json")
	check(err)
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&PO); err != nil {
		panic(err)
	}
	var blob []byte
	for i := 0; i < len(PO); i++ {
		//generating aks
		PO[i].Keys = &AKS{}
		PO[i].Keys.generate(db, path.Join("internal/role/PO/", PO[i].Id), PO[i].GetId())
		PO[i].Digest = base58.Encode(PO[i].GenerateDigest())
		// this will be done later
		//todo: generate signature
		blob = utils.GetAuthorize(PO[i].Digest, AuthorizeBy, db)
		PO[i].Signature = utils.BlobhashGeneration(blob)
		db.Put([]byte(PO[i].Signature), blob, nil)
	}
	return PO
}

func (m *VoterStruct) Verify(primaryKey string, db *leveldb.DB) bool {
	digest := sha256.New()
	for _, v := range m.Voters {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", m.Pscode)))
	digest.Write([]byte(fmt.Sprintf("%v", m.Seed)))

	return m.Signature == utils.BlobhashGeneration(utils.GetAuthorize(base58.Encode(digest.Sum(nil)), primaryKey, db))
}
func (m *PoBannedList) Verify(primaryKey string, db *leveldb.DB) bool {

	for i, v := range m.PoHash {
		if m.PoPoa[i] != utils.BlobhashGeneration(utils.GetAuthorize(v, primaryKey, db)) {
			return false
		}
	}
	return true
}

// todo: will be postquantum
// this interface take maximum of 3 parameter db, primarykey(opts), and location(opts)
func (v *AKS) generate(db *leveldb.DB, params ...string) {
	mode := schemes.ByName("Ed25519")
	public_key, private_key, err := mode.GenerateKey()
	if err != nil {
		panic(err)
	}
	//
	// Write the private key to a file named "private.key"
	rbin_private_key, err := private_key.MarshalBinary()
	if err != nil {
		panic(err)
	}
	rbin_public_key, err := public_key.MarshalBinary()
	if err != nil {
		panic(err)
	}
	v.PrivateKey = utils.BlobhashGeneration(rbin_private_key)
	v.PublicKey = utils.BlobhashGeneration(rbin_public_key)
	v.EncryptionAlgorithm = mode.Name()
	v.KeySize = int64(mode.PublicKeySize())
	v.SigSize = int64(mode.SignatureSize())

	// Write the private key to db
	err = db.Put([]byte(v.PrivateKey), rbin_private_key, nil)
	if err != nil {
		log.Panic(err)
	}
	// Write the public key to a db
	err = db.Put([]byte(v.PublicKey), rbin_public_key, nil)
	if err != nil {
		log.Panic(err)
	}

	if len(params) > 0 {
		if !utils.IsDirAvailable(params[0]) {
			err := os.Mkdir(params[0], os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		err = os.WriteFile(path.Join(params[0], "private.key"), rbin_private_key, 0644)
		if err != nil {
			log.Fatal(err)
		}
		// Write the public key to a file named "public.key"
		err = os.WriteFile(path.Join(params[0], "public.key"), rbin_public_key, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(params) > 1 {
		// Write the private key to db
		err = db.Put([]byte("privatekey_"+params[1]), rbin_private_key, nil)
		if err != nil {
			log.Panic(err)
		}
		// Write the public key to a db
		err = db.Put([]byte("publickey_"+params[1]), rbin_public_key, nil)
		if err != nil {
			log.Panic(err)
		}

	}
}

// todo: will be postquantum
func (v *AKS) Authorize(digest string, db *leveldb.DB) []byte {

	bit_privatekey, err := db.Get([]byte(v.PrivateKey), nil)
	if err != nil {
		panic(err)
	}
	mode := schemes.ByName("Ed25519")
	privatekey, err := mode.UnmarshalBinaryPrivateKey(bit_privatekey)
	if err != nil {
		panic(err)
	}
	// privatekey := custoPrivateKeyFromBytes(bit_privatekey)

	sig := mode.Sign(privatekey, []byte(digest), nil)
	if err != nil {
		panic(err)
	}
	//test sig

	return sig
}
