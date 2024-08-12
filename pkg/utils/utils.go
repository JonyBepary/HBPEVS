package utils

import (
	"crypto/sha256"
	"os"

	"github.com/cloudflare/circl/sign/schemes"
	"github.com/mr-tron/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

func String_hash_generation(str string) string {
	digest := sha256.New()
	digest.Write([]byte(str))
	return base58.Encode(digest.Sum(nil))
}

func BlobhashGeneration(bit []byte) string {
	hash := sha256.New()
	hash.Write(bit)
	return base58.Encode(hash.Sum(nil))
}

func GetAuthorize(digest, primaryKey string, db *leveldb.DB) []byte {

	bit_privatekey, err := db.Get([]byte(primaryKey), nil)
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

func IsDirAvailable(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

// ! need to be implemented
func Get_token(nid string) string {
	return base58.Encode([]byte(nid))
}

// Generate Digest of ReturningOfficer struct
// func roDigest(v ReturningOfficer) []byte {
// 	digest := sha256.New()
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Name)))
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Id)))
// 	for _, vr := range v.R_area {
// 		digest.Write([]byte(fmt.Sprintf("%v", vr)))
// 	}
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.PublicKey)))
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.EncryptionAlgorithm)))
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.KeySize)))
// 	digest.Write([]byte(fmt.Sprintf("%v", v.Keys.SigSize)))
// 	return digest.Sum(nil)
// }
