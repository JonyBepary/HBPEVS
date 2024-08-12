package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cloudflare/circl/sign/schemes"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenerateHeader(b *Block) BlockHeader {
	var Header BlockHeader
	Header.BlockID = 0
	Header.PreHash = "0"
	Header.Pscode = "8010"
	Header.Timestamp = time.Now().Unix()
	digest := sha256.New()

	digest.Write([]byte(fmt.Sprintf("%v", Header.BlockID)))
	digest.Write([]byte(fmt.Sprintf("%v", Header.PreHash)))
	digest.Write([]byte(fmt.Sprintf("%v", Header.Pscode)))

	//! timestamp caution and need to test throughly
	digest.Write([]byte(fmt.Sprintf("%v", Header.Timestamp)))
	digest.Write([]byte(fmt.Sprintf("%v", b.ElectionCommission.Digest)))

	for _, v := range b.ReturningOfficer {
		digest.Write([]byte(fmt.Sprintf("%v", v.Digest)))
	}
	for _, v := range b.PollingOfficer {
		digest.Write([]byte(fmt.Sprintf("%v", v.Digest)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", b.VoterStruct.Digest)))
	digest.Write([]byte(fmt.Sprintf("%v", b.CandidateStruct.Digest)))
	digest.Write([]byte(fmt.Sprintf("%v", b.ContractStruct.getDigest())))
	//! need to add PObanned list
	//! need to add Vote
	Header.PostHash = base58.Encode(digest.Sum(nil))
	return Header
}

func GenerateHeaderlb(lb *LiteBlock) BlockHeader {
	var Header BlockHeader
	Header.BlockID = 0
	Header.PreHash = "0"
	Header.Pscode = "8010"
	Header.Timestamp = time.Now().Unix()
	digest := sha256.New()

	digest.Write([]byte(fmt.Sprintf("%v", Header.BlockID)))
	digest.Write([]byte(fmt.Sprintf("%v", Header.PreHash)))
	digest.Write([]byte(fmt.Sprintf("%v", Header.Pscode)))

	//! timestamp caution and need to test throughly
	digest.Write([]byte(fmt.Sprintf("%v", Header.Timestamp)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.ElectionCommission)))

	for _, v := range lb.ReturningOfficer {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	for _, v := range lb.PollingOfficer {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", lb.VoterStruct)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.CandidateStruct)))

	digest.Write([]byte(fmt.Sprintf("%v", lb.ContractStruct)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.PoBannedList)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.Vote)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.SpendedToken)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.Token)))
	Header.PostHash = base58.Encode(digest.Sum(nil))

	return Header
}

// will be interfaced to CandidateStruct struct
func GenerateCandidate() CandidateStruct {
	var candidate CandidateStruct
	//parse info from configuration
	configFile, err := os.Open("config/candidate.json")
	check(err)
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&candidate); err != nil {
		panic(err)
	}
	digest := sha256.New()
	for _, v := range candidate.ApplicableCandidates {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	candidate.Digest = base58.Encode(digest.Sum(nil))
	return candidate
}

func (EC *ElectionCommission) GenerateElection(db *leveldb.DB) {
	//! need conv in post quantam
	mode := schemes.ByName("Ed25519")
	public_key, private_key, err := mode.GenerateKey()
	check(err)
	//
	if !utils.IsDirAvailable("internal/role/EC") {
		os.MkdirAll("internal/role/EC", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// Write the private key to a file named "private.key"
	rbin_private_key, err := private_key.MarshalBinary()
	check(err)
	rbin_public_key, err := public_key.MarshalBinary()
	check(err)
	// Write the private key to a file named "private.key"
	err = os.WriteFile("internal/role/EC/private.key", rbin_private_key, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// Write the public key to a file named "public.key"
	err = os.WriteFile("internal/role/EC/public.key", rbin_public_key, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Saving keys by their hash
	EC.PublicKey = utils.BlobhashGeneration(rbin_public_key)
	EC.PrivateKey = utils.BlobhashGeneration(rbin_private_key)
	EC.Ec_EncryptionAlgorithm = mode.Name()
	EC.KeySize = int64(mode.PublicKeySize())
	EC.SigSize = int64(mode.SignatureSize())

	// Saving keys by their hash
	err = db.Put([]byte(EC.PublicKey), rbin_public_key, nil)
	check(err)
	err = db.Put([]byte(EC.PrivateKey), rbin_private_key, nil)
	if err != nil {
		log.Panic(err)
	}

	// Saving hash of their for better access
	// like i'll ask for the hash using "EC_public_key"
	err = db.Put([]byte("EC_public_key"), rbin_public_key, nil)
	check(err)

	// like i'll ask for the hash using "EC_private_key"
	err = db.Put([]byte("EC_private_key"), rbin_private_key, nil)
	if err != nil {
		log.Panic(err)
	}

	// generating ec
	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", EC.PublicKey)))
	digest.Write([]byte(fmt.Sprintf("%v", EC.PrivateKey)))
	digest.Write([]byte(fmt.Sprintf("%v", EC.Ec_EncryptionAlgorithm)))
	digest.Write([]byte(fmt.Sprintf("%v", EC.KeySize)))
	digest.Write([]byte(fmt.Sprintf("%v", EC.SigSize)))
	EC.Digest = base58.Encode(digest.Sum(nil))
}

func parseSeed(port uint64, seed string, start string, duration string) (string, error) {

	baseURL := fmt.Sprintf("http://localhost:%d/", port)
	// seed := "def95e06826ffb028c97aa85096078e44c488e79e405626f2858a8b070c761c1"
	// start := "1681083456"
	// duration := "10800000000000"
	seed = strings.Replace(seed, " ", "", -1)
	start = strings.Replace(start, " ", "", -1)
	duration = strings.Replace(duration, " ", "", -1)

	url := (baseURL + "?seed=" + seed + "&start=" + start + "&duration=" + duration)

	// spew.Dump(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer res.Body.Close()
	// spew.Dump(res.Body)
	var response struct {
		Digest string `json:"digest"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {

		panic(err)
	}

	return response.Digest, nil
}
func (b *Block) getLiteBlock() *LiteBlock {
	//genrating light block
	var lb LiteBlock
	for i := 0; i < len(b.ReturningOfficer); i++ {
		lb.ReturningOfficer = append(lb.ReturningOfficer, b.ReturningOfficer[i].Digest)
	}
	for i := 0; i < len(b.PollingOfficer); i++ {
		lb.PollingOfficer = append(lb.PollingOfficer, b.PollingOfficer[i].Digest)
	}
	lb.Header = b.Header
	lb.ElectionCommission = b.ElectionCommission.Digest
	lb.ElectionTime = b.ElectionTime.Digest
	lb.CandidateStruct = b.CandidateStruct.Digest
	lb.VoterStruct = b.VoterStruct.Digest
	lb.SpendedToken = b.SpendedToken.Digest
	lb.PoBannedList = "PoBannedList"
	lb.ContractStruct = b.ContractStruct.getDigest()

	return &lb
}

func AddBlockAsLatest(b *Block, bindata *BinaryData) bool {
	if b == nil {
		return false
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := path.Join(wd, "BLOCK")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Working Directory: ", parent)
	db, err := leveldb.OpenFile(parent, nil)
	if err != nil {
		log.Printf("Failed to open database, Make sure your working directory is correct\n")
		panic(err)
	}
	defer db.Close()
	blob, err := db.Get([]byte("latest_lite_block"), nil)
	if err != nil {
		panic(err)
	}
	var tmp_lb LiteBlock

	tmp_b := tmp_lb.GetFullBlock()
	if err != nil {
		return false
	}

	if tmp_b.Header.BlockID >= b.Header.BlockID {
		log.Println("Block ID is not greater than latest block")
		return false
	}
	if tmp_b.Header.BlockID+1 != b.Header.BlockID {
		log.Println("Block ID is not greater than latest block")
		return false
	}
	if tmp_b.Header.PostHash != b.Header.PreHash {
		log.Println("Block hash is not equal to previous block hash")
		return false
	}
	if tmp_b.Header.ElectionID != b.Header.ElectionID {
		log.Println("Block Election ID is not equal to previous block Election ID")
		return false
	}
	if tmp_b.Header.Pscode != b.Header.Pscode {
		log.Println("Block Pscode is not equal to previous block Pscode")
		return false
	}
	if b.Vote == "" {
		log.Println("Block Vote is nil")
		return false
	}
	if b.Token == "" {
		log.Println("Block Token is nil")
		return false
	}
	tmp_b.Header = b.Header

	if b.CandidateStruct != nil {
		tmp_b.CandidateStruct = b.CandidateStruct
		err = db.Put([]byte(tmp_b.CandidateStruct.Signature), blob, nil)
		if err != nil {
			log.Panic(err)
		}
		err = db.Put([]byte("candidate_signature"), blob, nil)
		if err != nil {
			log.Panic(err)
		}

	}
	if b.ContractStruct != nil {
		tmp_b.ContractStruct = b.ContractStruct
		for i := 0; i < len(tmp_b.ContractStruct.ProofOfCommand); i++ {
			err := db.Put([]byte(tmp_b.ContractStruct.ProofOfCommand[i].ProofOfCommand), bindata.ContractProofOfCommand[i], nil)
			if err != nil {
				panic(err)
			}
			err = db.Put([]byte("proof_of_command_"+fmt.Sprint(i+1)), blob, nil)
			if err != nil {
				log.Panic(err)
			}
		}
	}

	if b.PoBannedList != nil {
		tmp_b.PoBannedList = b.PoBannedList
		for i := 0; i < len(tmp_b.PoBannedList.GetPoHash()); i++ {
			err := db.Put([]byte(tmp_b.PoBannedList.PoHash[i]), bindata.PoBannedList[i], nil)
			if err != nil {
				log.Panic(err)
			}
		}
		blob, err = tmp_b.PoBannedList.XXX_Marshal(nil, true)
		if err != nil {
			log.Panic(err)
		}
		err = db.Put([]byte("PoBannedList"), blob, nil)
		if err != nil {
			log.Panic(err)
		}
	}
	if b.PollingOfficer != nil {
		tmp_b.PollingOfficer = b.PollingOfficer
		for i := 0; i < len(tmp_b.PollingOfficer); i++ {

			err := db.Put([]byte(tmp_b.PollingOfficer[i].Signature), bindata.PollingOfficerSignature[i], nil)
			if err != nil {
				log.Panic(err)
			}
			err = db.Put([]byte("signature_"+tmp_b.PollingOfficer[i].Id), blob, nil)
			if err != nil {
				log.Panic(err)
			}
			blob, err = tmp_b.PollingOfficer[i].XXX_Marshal(nil, true)
			check(err)
			err = db.Put([]byte(fmt.Sprint(tmp_b.PollingOfficer[i].Digest)), blob, nil)
			check(err)
			err = db.Put([]byte(fmt.Sprint(tmp_b.PollingOfficer[i].Id)), blob, nil)
			check(err)
		}

	}
	
	tmp_b.Token = b.Token
	tmp_b.Vote = b.Vote

	blob, err = db.Get([]byte(tmp_b.SpendedToken.Token), nil)
	if err != nil {
		log.Panic(err)
	}
	var spended_token treeset.Set
	err = spended_token.UnmarshalJSON(blob)
	if err != nil {
		log.Panic(err)
	}
	spended_token.Add(tmp_b.Token)

	tmp_b.SpendedToken.AddItemWithDigestProofUpdate(tmp_b.Token, tmp_b.ReturningOfficer[0].Id, db)

	blob, err = tmp_b.SpendedToken.XXX_Marshal(nil, true)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("SpendedToken"), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	blob, err = tmp_b.XXX_Marshal(nil, true)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("latest"), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte(fmt.Sprintf("block_%v", tmp_b.Header.BlockID)), []byte(fmt.Sprint(tmp_b.Header.BlockID)), nil)
	check(err)
	lb := tmp_b.getLiteBlock()
	blob, err = lb.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(fmt.Sprintf("lite_block_%d", lb.Header.BlockID)), blob, nil)
	check(err)
	err = db.Put([]byte("latest_lite_block"), blob, nil)
	check(err)
	return true
}
