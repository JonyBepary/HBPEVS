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

// GenerateToken will generate a token for the VOTER
// It'll take the hash of the voter and generate a token
// which will be used to verify the voter
func GenerateToken(digest, start, duration string) string {
	// var psig, rsig, token string

	psig, err := parseSeed(8112, digest, start, duration)
	if err != nil {
		panic(err)
	}
	rsig, err := parseSeed(8113, psig, start, duration)
	if err != nil {
		panic(err)
	}
	token, err := parseSeed(8114, rsig, start, duration)
	if err != nil {
		panic(err)
	}

	return token
}

func (v *SpendedToken) AppendDigest() []byte {

	digest := sha256.New()
	digest.Write([]byte(v.Digest))
	digest.Write([]byte(fmt.Sprintf("%v", v.Token[len(v.Token)-1])))
	return digest.Sum(nil)
}

func (v *SpendedToken) GenerateDigestALL(db *leveldb.DB) []byte {

	digest := sha256.New()
	tokenSet := v.GetTree(db)
	for _, v := range tokenSet.Values() {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}

	return digest.Sum(nil)
}

func (v *SpendedToken) UpdateDigest(db *leveldb.DB) {

	digest := sha256.New()
	tokenSet := v.GetTree(db)
	for _, v := range tokenSet.Values() {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	v.Digest = base58.Encode(digest.Sum(nil))
}

func (v *SpendedToken) AddItemWithDigestProofUpdate(str, authorizeBy string, db *leveldb.DB) {
	tree := v.GetTree(db)
	tree.Add(str)
	tkn, err := tree.ToJSON()
	if err != nil {
		panic(err)
	}
	v.Token = utils.BlobhashGeneration(tkn)

	err = db.Put([]byte(v.Token), tkn, nil)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("spended_token"), tkn, nil)
	if err != nil {
		panic(err)
	}

	// v.UpdateDigest(db)
	blob := utils.GetAuthorize(v.Token, authorizeBy, db)
	v.TokenProof = utils.BlobhashGeneration(blob)
	err = db.Put([]byte(v.TokenProof), blob, nil)
	if err != nil {
		panic(err)
	}
	err = db.Put([]byte("spended_token_proof"), blob, nil)

	if err != nil {
		panic(err)
	}
	digest := sha256.New()
	digest.Write([]byte(v.Digest))
	digest.Write([]byte(v.Token))
	v.Digest = base58.Encode(digest.Sum(nil))

}

func (v *SpendedToken) GetTree(db *leveldb.DB) *treeset.Set {

	blob, err := db.Get([]byte(v.Token), nil)
	if err != nil {
		log.Panic(err)
	}
	tree := treeset.NewWithStringComparator()
	if err := tree.UnmarshalJSON(blob); err != nil {
		panic(err)
	}
	return tree
}
