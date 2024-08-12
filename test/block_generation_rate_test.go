package main

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/dchest/uniuri"
	"github.com/syndtr/goleveldb/leveldb"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var Candidate = []string{"006", "007", "013", "002", "003", "011", "008", "014"}

type VoterStruct1 struct {
	NID      string
	Name     string
	Filename string
}

var v_list *list.List

func make_struct(n uint64, pscode string, db *leveldb.DB) {
	var i uint64 = 0
	var base uint64 = 20215103000

	for ; i < n; i++ {

		// fmt.Println(i)
		rand.Seed(time.Now().UnixNano())
		name := uniuri.NewLenChars(12, StdChars)
		nid := strconv.Itoa(int(base + i))

		hash := sha256.New()
		hash.Write([]byte(fmt.Sprint(nid)))
		hash.Write([]byte(fmt.Sprint(pscode)))
		filename := fmt.Sprintf("%x", hash.Sum(nil))

		Bbytes, err := json.Marshal(VoterStruct1{NID: nid, Name: name, Filename: filename})
		if err != nil {
			panic(err)
		}
		db.Put([]byte(fmt.Sprintf("%d", i)), Bbytes, nil)
	}

}

func test_sword_of_durant(nid, name, pscode, PO, union, thana, district string) {

	url := fmt.Sprintf("http://localhost:8888/sword_of_durant?nid=%v&name=%v&pscode=%v&PO=%v&union=%v&thana=%v&district=%v", nid, name, pscode, PO, union, thana, district)
	// fmt.Println(url)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/home/jony/Pictures/pic.png")
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("avatar", filepath.Base("/home/jony/Pictures/nid_pic/rony.jpg"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(body)
		fmt.Println(err)
		return
	}

}

func Benchmark_NID_SERVER(b *testing.B) {

	pscode := "8010"
	PO := "Naodoba"
	union := "Naodoba"
	thana := "Janjira"
	district := "Shariatpur"
	var n uint64 = uint64(b.N)
	db, err := leveldb.OpenFile("levelDB", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	make_struct(n, pscode, db)
	for i := uint64(0); i < n; i++ {
		data, err := db.Get([]byte(fmt.Sprintf("%d", i)), nil)
		if err != nil {
			fmt.Println(err)
		}
		v := VoterStruct1{}
		err = json.Unmarshal(data, &v)
		if err != nil {
			fmt.Println(err)
		}
		test_sword_of_durant(v.NID, v.Name, pscode, PO, union, thana, district)
	}

}

func test_addVoter(vote, pscode, digest, po_id, po_privkey string) {

	url := "http://localhost:5669/add_vote?vote=" + vote + "&pscode=" + pscode + "&digest=" + digest + "&po_id=" + po_id + "&po_privkey=" + po_privkey + "&verbose=0"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func Benchmark_BlockGenearationRate(b *testing.B) {

	vote := "006"
	pscode := "8010"
	po_id := "PO123"
	po_privkey := "privatekey_PO123"
	/*
		go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> nid_data.txt
		go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> newblock_data.txt -timeout 60m
		curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681887262&duration=10800&config=2&verbose=2'; echo ""

	*/
	//po_pubkey: ""
	//posig:MGYCMQDcJsixoA/8g2JxizoyyEXZwrk49tY8a1I/0I6dzwkWfUbl2vlNPOiJb6oJ6ul1mFcCMQD9AB7U8TmkNuZktssDkHdfwZltwODbcuuDAa/FCXeROoTltaDF5DPi6ucIcgyevlc=
	// pscode := "8010"
	// PO := "Naodoba"
	// union := "Naodoba"
	// thana := "Janjira"
	// district := "Shariatpur"

	var n uint64 = uint64(b.N)

	db, err := leveldb.OpenFile("levelDB", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	for i := uint64(0); i < n; i++ {
		data, err := db.Get([]byte(fmt.Sprintf("%d", i)), nil)
		if err != nil {
			fmt.Println(err)
		}
		v := VoterStruct1{}
		err = json.Unmarshal(data, &v)
		if err != nil {
			fmt.Println(err)
		}

		test_addVoter(vote, pscode, v.Filename, po_id, po_privkey)
	}

}

// func TestAddBlock(t *testing.T) {
// 	// Create a temporary test database directory
// 	dir, err := ioutil.TempDir("", "test_db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.RemoveAll(dir)

// 	// Initialize a new blockchain
// 	chain := NewBlockchain(dir)

// 	// Create a new block
// 	block := &Block{
// 		Header: BlockHeader{
// 			BlockID:        1,
// 			Timestamp:      1234567890,
// 			Pre_hash:       "",
// 			Post_hash:      "",
// 			ProofOfBlock: "",
// 			Pscode:         "123",
// 		},
// 		Candidate:        CandidateStruct{},
// 		EC:               ElectionCommission{},
// 		ReturningOfficer: []ReturningOfficer{},
// 		PollingOfficer:  []PollingOfficer{},
// 		VoterList:       VoterStruct{},
// 		Contract:         ContractStruct{},
// 		PoBannedList:     PoBannedList{},
// 		Token:            "",
// 		Vote:             "",
// 	}

// 	// Add the block to the blockchain
// 	err = chain.AddBlock(block)
// 	if err != nil {
// 		t.Errorf("AddBlock error: %v", err)
// 	}

// 	// Get the last block from the blockchain
// 	lastBlock, err := chain.GetLastBlock()
// 	if err != nil {
// 		t.Errorf("GetLastBlock error: %v", err)
// 	}

// 	// Check if the last block added to the blockchain is the same as the block we added
// 	if lastBlock.Header.BlockID != block.Header.BlockID {
// 		t.Errorf("Expected block ID %d, got %d", block.Header.BlockID, lastBlock.Header.BlockID)
// 	}
// 	if lastBlock.Header.Timestamp != block.Header.Timestamp {
// 		t.Errorf("Expected timestamp %d, got %d", block.Header.Timestamp, lastBlock.Header.Timestamp)
// 	}
// 	if lastBlock.Header.PreHash != block.Header.PreHash {
// 		t.Errorf("Expected pre-hash %s, got %s", block.Header.PreHash, lastBlock.Header.PreHash)
// 	}
// 	if lastBlock.Header.PostHash != block.Header.PostHash {
// 		t.Errorf("Expected post-hash %s, got %s", block.Header.PostHash, lastBlock.Header.PostHash)
// 	}
// 	if lastBlock.Header.ProofOfBlock != block.Header.ProofOfBlock {
// 		t.Errorf("Expected proof of block %s, got %s", block.Header.ProofOfBlock, lastBlock.Header.ProofOfBlock)
// 	}
// 	if lastBlock.Header.Pscode != block.Header.Pscode {
// 		t.Errorf("Expected pscode %s, got %s", block.Header.Pscode, lastBlock.Header.Pscode)
// 	}

// 	// Get the block from the blockchain by block ID
// 	blockByID, err := chain.GetBlockByID(block.Header.BlockID)
// 	if err != nil {
// 		t.Errorf("GetBlockByID error: %v", err)
// 	}

// 	// Check if the block we added is the same as the block we got by ID
// 	if !jsonEqual(block, blockByID) {
// 		t.Errorf("Expected block %v, got %v", block, blockByID)
// 	}
// }

// // jsonEqual checks whether two JSON-encoded values are equal.
// func jsonEqual(a, b interface{}) bool {
// 	// Marshal the values to JSON.
// 	aJSON, err := json.Marshal(a)
// 	if err != nil {
// 		return false
// 	}
// 	bJSON, err := json.Marshal(b)
// 	if err != nil {
// 		return false
// 	}

// 	// Unmarshal the JSON-encoded values into generic interface{}.
// 	var aIntf, bIntf interface{}
// 	err = json.Unmarshal(aJSON, &aIntf)
// 	if err != nil {
// 		return false
// 	}
// 	err = json.Unmarshal(bJSON, &bIntf)
// 	if err != nil {
// 		return false
// 	}

// 	// Compare the unmarshalled values for equality.
// 	return reflect.DeepEqual(aIntf, bIntf)
// }
