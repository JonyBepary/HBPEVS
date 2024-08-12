package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/internal/blockchain"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var voterSet *treeset.Set
var tokenSet *treeset.Set

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fast_add_vote(c *gin.Context) {

	po_id := c.Query("po_id")
	fmt.Println("po_id: " + po_id)

	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}
	defer db.Close()

	bits, err := db.Get([]byte("latest"), nil)
	if err != nil {
		panic(err)
	}

	var lb blockchain.LiteBlock
	err = lb.XXX_Unmarshal(bits)
	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}

	bits, err = db.Get([]byte(fmt.Sprint(lb.ElectionTime)), nil)
	if err != nil {
		panic(err)
	}

	var election_time blockchain.ElectionTime
	err = election_time.XXX_Unmarshal(bits)
	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}

	voter := c.Query("digest")
	lb.Header.Pscode = c.Query("pscode")
	lb.Vote = c.Query("vote")
	lb.Header.PreHash = lb.Header.PostHash
	lb.Header.BlockID += 1
	lb.Header.Timestamp = time.Now().Unix()
	PO_PrivateKey := c.Query("po_privkey")
	// getting primary key of 'poling officer privatekey'
	PO_obj_blob, err := db.Get([]byte(po_id), nil)
	if err != nil {
		panic(err)
	}
	var PO_obj blockchain.PollingOfficer
	err = PO_obj.XXX_Unmarshal(PO_obj_blob)
	if err != nil {
		panic(err)
	}

	// get this pro Authorize by PrivateKey of a polling officer
	// fmt.Println("lb: " + lb.Header.PostHash)
	// fmt.Println("PrivateKey: " + PrivateKey)

	// lbBytes, err := db.Get([]byte(lb.VoterStruct), nil)
	// var voters blockchain.VoterStruct
	// voters.XXX_Unmarshal(lbBytes)
	// is_valid := false
	// for _, v := range voters.Voters {
	// 	if v == voter {
	// 		is_valid = true
	// 	}
	// }
	// checking if voter is valid or not
	if !voterSet.Contains(voter) {
		c.String(500, "Voter is not valid")
		return
	}

	lb.Token = blockchain.GenerateToken(voter, election_time.Start, election_time.Duration)

	lbBytes, err := db.Get([]byte(fmt.Sprintf("spendedtoken")), nil)
	if err != nil {
		panic(err)
	}
	var spended_token blockchain.SpendedToken
	err = spended_token.XXX_Unmarshal(lbBytes)
	if err != nil {
		panic(err)
	}

	// checking if token is valid or not
	if condition := tokenSet.Contains(lb.Token); condition {
		c.String(500, "Voter already voted")
		return
	}

	// checking if token already used or not
	// for _, v := range spended_token.Token {
	// 	if v == lb.Token {
	// 		c.String(500, fmt.Sprintf("Voter already voted"))
	// 		return
	// 	}
	// }

	// adding token to spended token
	tokenSet.Add(lb.Token)

	spended_token.Token = append(spended_token.Token, lb.Token)
	spended_token.Digest = base58.Encode(spended_token.AppendDigest())
	spended_tokenBytes, err := spended_token.XXX_Marshal(nil, true)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte(spended_token.Digest), spended_tokenBytes, nil)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte(fmt.Sprintf("spendedtoken")), spended_tokenBytes, nil)
	if err != nil {
		panic(err)
	}

	// posthash
	lb.Header.PostHash = blockchain.Generate_block_hash(lb)

	// proof of block generation
	blob := utils.GetAuthorize(lb.Header.PostHash, PO_PrivateKey, db)
	lb.Header.ProofOfBlock = utils.BlobhashGeneration(blob)
	lb.SpendedToken = spended_token.Digest
	lb.Checkblock(blob)

	err = db.Put([]byte(lb.Header.PostHash), blob, nil)
	if err != nil {
		panic(err)
	}

	lbBytes, err = lb.XXX_Marshal(nil, true)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte(fmt.Sprintf("block_%v", lb.Header.BlockID)), lbBytes, nil)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("latest"), lbBytes, nil)
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}
	switch c.Query("verbose") {
	case "0":
	case "1":
		c.String(http.StatusOK, spew.Sdump(lb))

	default:
		c.String(http.StatusOK, "Vote added successfully")
	}
}

func add_vote(c *gin.Context) {

	po_id := c.Query("po_id")
	fmt.Println("po_id: " + po_id)

	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}
	defer db.Close()
	bits, err := db.Get([]byte("latest"), nil)
	if err != nil {
		panic(err)
	}
	var lb blockchain.LiteBlock
	if err != nil {
		panic(err)
	}
	err = lb.XXX_Unmarshal(bits)
	if err != nil {
		panic(err)
	}

	voter := c.Query("digest")
	lb.Header.Pscode = c.Query("pscode")
	lb.Vote = c.Query("vote")
	lb.Header.PreHash = lb.Header.PostHash
	lb.Header.BlockID += 1

	bits, err = db.Get([]byte(fmt.Sprint(lb.ElectionTime)), nil)
	if err != nil {
		panic(err)
	}
	var election_time blockchain.ElectionTime
	err = election_time.XXX_Unmarshal(bits)

	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}
	// spew.Dump(lb)

	// getting primary key of 'poling officer privatekey'
	po_private_key := c.Query("po_privkey")

	PO_obj_blob, err := db.Get([]byte(po_id), nil)
	if err != nil {
		panic(err)
	}
	var PO_obj blockchain.PollingOfficer
	err = PO_obj.XXX_Unmarshal(PO_obj_blob)
	if err != nil {
		panic(err)
	}

	// get this pro Authorize by po_private_key of a polling officer
	// fmt.Println("lb: " + lb.Header.PostHash)
	// fmt.Println("po_private_key: " + po_private_key)

	// lbBytes, err := db.Get([]byte(lb.VoterStruct), nil)
	// var voters blockchain.VoterStruct
	// voters.XXX_Unmarshal(lbBytes)
	// is_valid := false
	// for _, v := range voters.Voters {
	// 	if v == voter {
	// 		is_valid = true
	// 	}
	// }
	// checking if voter is valid or not
	if !voterSet.Contains(voter) {
		c.String(500, "Voter is not valid")
		return
	}

	lb.Token = blockchain.GenerateToken(voter, election_time.Start, election_time.Duration)

	// checking if token is valid or not
	if condition := tokenSet.Contains(lb.Token); condition {
		c.String(500, "Voter already voted")
		return
	}

	lb.Header.Timestamp = time.Now().Unix()
	lb.Header.PostHash = blockchain.Generate_block_hash(lb)

	// lbBytes, err = db.Get([]byte(fmt.Sprintf("spendedtoken")), nil)
	// if err != nil {
	// 	panic(err)
	// }

	// var spended_token blockchain.SpendedToken
	// spended_token.XXX_Unmarshal(lbBytes)
	// checking if token already used or not
	// for _, v := range spended_token.Token {
	// 	if v == lb.Token {
	// 		c.String(500, fmt.Sprintf("Voter already voted"))
	// 		return
	// 	}
	// }

	// adding token to spended token
	// spended_token.Token = append(spended_token.Token, lb.Token)
	// tokenSet.Add(lb.Token)
	// spended_tokenBytes, err := spended_token.XXX_Marshal(nil, true)
	// if err != nil {
	// 	panic(err)
	// }
	// pktokens := base58.Encode([]byte(spended_token.GenerateDigest()))

	// err = db.Put([]byte(fmt.Sprintf(pktokens)), spended_tokenBytes, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// err = db.Put([]byte(fmt.Sprintf("spendedtoken")), spended_tokenBytes, nil)
	// if err != nil {
	// 	panic(err)
	// }

	blob := utils.GetAuthorize(lb.Header.PostHash, po_private_key, db)
	// posthash
	lb.Header.ProofOfBlock = utils.BlobhashGeneration(blob)

	lb.Checkblock(blob)

	err = db.Put([]byte(lb.Header.PostHash), blob, nil)
	if err != nil {
		panic(err)
	}

	lbBytes, err := lb.XXX_Marshal(nil, true)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte(fmt.Sprintf("block_%v", lb.Header.BlockID)), lbBytes, nil)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("latest"), lbBytes, nil)
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
	switch c.Query("verbose") {
	case "0":
	case "1":
		c.String(http.StatusOK, spew.Sdump(lb))

	default:
		c.String(http.StatusOK, "Vote added successfully")
	}
}
func initBlock(c *gin.Context) {
	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	pscode := c.Query("pscode")
	seed := c.Query("seed")
	config := c.Query("config")
	start := c.Query("start")
	duration := c.Query("duration")
	block := new(blockchain.Block)
	block.CreateGenesisBlock(start, duration, config, pscode, seed, db)
	for _, v := range block.VoterStruct.Voters {
		voterSet.Add(v)
	}
	if err != nil {
		panic(err)
	}
	switch c.Query("verbose") {
	case "1":
		res, err := json.Marshal(block)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, string(res))

	case "2":
		c.String(http.StatusOK, block.PollingOfficer[0].Keys.PrivateKey)
	default:
		c.String(http.StatusOK, "Genesis Block Created")
	}

}
func Test_initBlock(c *gin.Context) {
	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	pscode := c.Query("pscode")
	seed := c.Query("seed")
	start := c.Query("start")
	duration := c.Query("duration")
	config := c.Query("config")
	block := new(blockchain.Block)

	block.CreateGenesisBlock(start, duration, config, pscode, seed, db)
	// c.Writer.Write([]byte())

}

func main() {
	// The returned DB instance is safe for concurrent use. Which mean that all
	// DB's methods may be called concurrently from multiple goroutine.
	// Test_initBlock()
	voterSet = treeset.NewWithStringComparator()
	tokenSet = treeset.NewWithStringComparator()

	router := gin.Default()
	router.POST("/add_vote", fast_add_vote)
	router.POST("/generate_genesis", initBlock)
	router.Run(":5669")
	// add_vote(0)
	// filename := "./genesis.json"
	// os.WriteFile(filename, file, 0644)

}

//todo: create_election: This API function allows the polling officer to create a new election. The function takes as input the details of the election, such as the title, start and end dates, and the list of candidates.

//todo: cast_vote: This API function allows the voter to cast their vote. The function takes as input the voter's ID and the ID of the candidate they are voting for.

//todo: end_election: This API function ends the election and stops the voting process. The function takes no input parameters.

//todo: tally_votes: This API function tallies the votes and generates the election results. The function takes no input parameters and returns the final results of the election.

//? todo: verify_voter: This API function verifies the eligibility of a voter to participate in the election. The function takes as input the voter's ID and other relevant details, and checks that they meet the eligibility criteria.

//? todo: get_voter_status: This API function retrieves the current status of a voter, such as whether they have already cast their vote or not. The function takes as input the voter's ID and returns their current status.

//? todo: cancel_vote: This API function allows a voter to cancel their vote if they made a mistake or changed their mind. The function takes as input the voter's ID and the ID of the candidate they originally voted for.

//? todo: audit_election: This API function provides an auditable trail of the entire election process, including the details of all voters, candidates, and votes. The function takes no input parameters and returns a complete audit log of the election.
