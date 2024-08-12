package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/gin-gonic/gin"

	"github.com/sohelahmedjoni/pqhbevs_hac/internal/blockchain"
	"github.com/sohelahmedjoni/pqhbevs_hac/internal/network"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
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

	blob, err := db.Get([]byte("latest_lite_block"), nil)
	if err != nil {
		panic(err)
	}

	var lb blockchain.LiteBlock
	err = lb.XXX_Unmarshal(blob)
	if err != nil {
		c.String(500, err.Error())
		panic(err)
	}

	blob, err = db.Get([]byte(fmt.Sprint(lb.ElectionTime)), nil)
	if err != nil {
		panic(err)
	}

	var election_time blockchain.ElectionTime
	err = election_time.XXX_Unmarshal(blob)
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
	// PO_obj_blob, err := db.Get([]byte(po_id), nil)
	// if err != nil {
	// 	panic(err)
	// }
	// var PO_obj blockchain.PollingOfficer
	// err = PO_obj.XXX_Unmarshal(PO_obj_blob)
	// if err != nil {
	// 	panic(err)
	// }

	// var votersStruct blockchain.VoterStruct
	// blob, err = db.Get([]byte(lb.VoterStruct), nil)
	// votersStruct.XXX_Unmarshal(blob)

	// checking if voter is valid or not
	if !voterSet.Contains(voter) {
		c.String(500, "Voter is not valid")
		return
	}
	lb.Token = blockchain.GenerateToken(voter, election_time.Start, election_time.Duration)
	// checking if token is valid or not
	lbBytes, err := db.Get([]byte(lb.SpendedToken), nil)
	if err != nil {
		panic(err)
	}
	var spended_token blockchain.SpendedToken
	err = spended_token.XXX_Unmarshal(lbBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("spended_token: \n")
	blob, err = db.Get([]byte(spended_token.Token), nil)
	if err != nil {
		panic(err)
	}
	tr := treeset.NewWithStringComparator()
	err = tr.FromJSON(blob)
	if err != nil {
		panic(err)
	}

	// checking if token is valid or not
	if condition := tokenSet.Contains(lb.Token); condition {
		c.String(500, "Voter already voted")
		return
	}
	//todo: ro selection from query
	blob, err = db.Get([]byte(lb.ReturningOfficer[0]), nil)
	if err != nil {
		panic(err)
	}
	var ro blockchain.ReturningOfficer
	err = ro.XXX_Unmarshal(blob)
	if err != nil {
		panic(err)
	}

	spended_token.AddItemWithDigestProofUpdate(lb.Token, ro.Id, db)

	// adding token to spended token
	spended_tokenBytes, err := spended_token.XXX_Marshal(nil, true)
	if err != nil {
		panic(err)
	}
	lb.SpendedToken = spended_token.Digest
	err = db.Put([]byte(lb.SpendedToken), spended_tokenBytes, nil)
	if err != nil {
		panic(err)
	}
	err = db.Put([]byte("spendedtoken"), spended_tokenBytes, nil)
	if err != nil {
		panic(err)
	}

	// posthash generation and proof of block generation
	lb.Header.PostHash = blockchain.Generate_block_hash(lb)

	// proof of block generation
	blob = utils.GetAuthorize(lb.Header.PostHash, PO_PrivateKey, db)
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

	err = db.Put([]byte(fmt.Sprintf("liteblock_%v", lb.Header.BlockID)), lbBytes, nil)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("latest_lite_block"), lbBytes, nil)
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
		c.String(200, spew.Sdump(lb))

	default:
		c.String(200, "Vote added successfully")
	}
}

/*
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

	lbBytes, err := db.Get([]byte(fmt.Sprintf(lb.SpendedToken)), nil)
	check(err)
	err = tokenSet.UnmarshalJSON(lbBytes)
	check(err)
	lb.Token = blockchain.GenerateToken(voter, election_time.Start, election_time.Duration)

	// checking if token is valid or not
	if condition := tokenSet.Contains(lb.Token); condition {
		c.String(500, "Voter already voted")
		return
	}
	tokenSet.Add(lb.Token)
	blob, err := db.Get([]byte(lb.SpendedToken), nil)
	check(err)
	spended_token := blockchain.SpendedToken{}
	spended_token.XXX_Unmarshal(blob)
	spended_token.AddItemWithDigestProofUpdate(lb.Token, po_id , db)
	err = db.Put([]byte(fmt.Sprintf(lb.SpendedToken)), blob, nil)
	check(err)
	lb.Header.Timestamp = time.Now().Unix()
	lb.Header.PostHash = blockchain.Generate_block_hash(lb)

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
		c.String(200, spew.Sdump(lb))

	default:
		c.String(200, "Vote added successfully")
	}
}
*/
// initBlock is used to initialize the block
// it will create genesis block

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
		c.String(200, string(res))

	case "2":
		c.String(200, block.PollingOfficer[0].Keys.PrivateKey)
	default:
		c.String(200, "Genesis Block Created")
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
}

func getfullblock() {
	// start p2p network
	// The returned DB instance is safe for concurrent use
	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	bit, err := db.Get([]byte("latest"), nil)
	check(err)
	block := new(blockchain.LiteBlock)
	block.XXX_Unmarshal(bit)

	fullblock := new(blockchain.Block)
	fullblock.Header = block.Header

	// ElectionTime get
	bit, err = db.Get([]byte(block.ElectionTime), nil)
	check(err)
	fullblock.ElectionTime = new(blockchain.ElectionTime)
	fullblock.ElectionTime.XXX_Unmarshal(bit)

	// CANDIDATE get
	bit, err = db.Get([]byte(block.CandidateStruct), nil)
	check(err)
	fullblock.CandidateStruct = new(blockchain.CandidateStruct)
	fullblock.CandidateStruct.XXX_Unmarshal(bit)

	// VOTER get
	bit, err = db.Get([]byte(block.VoterStruct), nil)
	check(err)
	fullblock.VoterStruct = new(blockchain.VoterStruct)
	fullblock.VoterStruct.XXX_Unmarshal(bit)

	// POLLINGOFFICER get
	fullblock.PollingOfficer = make([]*blockchain.PollingOfficer, len(block.PollingOfficer))
	for i := 0; i < len(block.PollingOfficer); i++ {
		bit, err = db.Get([]byte(block.PollingOfficer[i]), nil)
		check(err)
		fullblock.PollingOfficer[i] = new(blockchain.PollingOfficer)
		fullblock.PollingOfficer[i].XXX_Unmarshal(bit)
	}

	// returning officer get
	fullblock.ReturningOfficer = make([]*blockchain.ReturningOfficer, len(block.ReturningOfficer))
	for i := 0; i < len(block.ReturningOfficer); i++ {
		bit, err = db.Get([]byte(block.ReturningOfficer[i]), nil)
		check(err)
		fullblock.ReturningOfficer[i] = new(blockchain.ReturningOfficer)
		fullblock.ReturningOfficer[i].XXX_Unmarshal(bit)
	}

	// get po banned list
	bit, err = db.Get([]byte(block.PoBannedList), nil)
	check(err)
	fullblock.PoBannedList = new(blockchain.PoBannedList)
	fullblock.PoBannedList.XXX_Unmarshal(bit)

	// get spent token
	bit, err = db.Get([]byte(block.SpendedToken), nil)
	check(err)
	fullblock.SpendedToken = new(blockchain.SpendedToken)
	fullblock.SpendedToken.XXX_Unmarshal(bit)
	// get token
	fullblock.Token = block.Token
	// get token
	fullblock.Vote = block.Vote

	b, err := json.MarshalIndent(fullblock, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}

func main() {
	// start p2p network
	// cfg := network.Config{
	// 	RendezvousString: "meetme",
	// 	ProtocolID:       "/block/blockreq/0.0.1",
	// 	ListenHost:       "0,0,0,0",
	// 	ListenPort:       4001,
	// 	ApiPort:          5669,

	// }
	// go network.NodeInstance(cfg)

	help := flag.Bool("help", false, "Display Help")
	rnd := flag.Bool("rnd", false, "start node with a random free port")
	nop2p := flag.Bool("nop2p", false, "start node without p2p")
	cfg := parseFlags()
	cfg.RendezvousString = ""
	if *help {
		fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
		fmt.Printf("Usage: \n   Run './app'\nor Run './app -port [port] -rendezvous [string] -pid [protocol ID] or just ./app -rnd '\n")

		os.Exit(0)
	}
	if *rnd {
		rnd := rand.New(rand.NewSource(666))
		// Choose random ports between 10000-10100
		cfg.ListenPort = rnd.Intn(5000) + 10000
		cfg.ApiPort = rnd.Intn(5000) + 4999
	}
	// start p2p network

	// The returned DB instance is safe for concurrent use. Which mean that all
	// DB's methods may be called concurrently from multiple goroutine.
	// Test_initBlock()
	voterSet = treeset.NewWithStringComparator()
	tokenSet = treeset.NewWithStringComparator()

	if utils.IsDirAvailable("BLOCK") {
		db, err := leveldb.OpenFile("BLOCK", nil)
		check(err)

		bit, err := db.Get([]byte("latest_lite_block"), nil)
		check(err)
		block := new(blockchain.LiteBlock)
		block.XXX_Unmarshal(bit)
		//! spew.Dump(block)
		if block.Header.BlockID < 1 {
			fmt.Println("no genesis block")
			*nop2p = true
		} else {
			fmt.Println("genesis block exist")
			blob, err := db.Get([]byte(block.SpendedToken), nil)
			check(err)

			spendedToken := new(blockchain.SpendedToken)
			spendedToken.XXX_Unmarshal(blob)

			blob, err = db.Get([]byte(spendedToken.Token), nil)
			check(err)
			tokenSet.FromJSON(blob)

			blob, err = db.Get([]byte(block.VoterStruct), nil)
			check(err)
			voterStruct := new(blockchain.VoterStruct)
			voterStruct.XXX_Unmarshal(blob)

			for i := 0; i < len(voterStruct.Voters); i++ {
				voterSet.Add(voterStruct.Voters[i])
			}
		}
		db.Close()
	} else {
		*nop2p = true
	}

	if !*nop2p {

		go network.NodeInstance(cfg)
	}
	router := gin.Default()
	router.POST("/add_vote", fast_add_vote)
	router.POST("/generate_genesis", initBlock)
	router.Run(fmt.Sprintf(":%d", cfg.ApiPort))

	// // // add_vote(0)
	// // // filename := "./genesis.json"
	// // // os.WriteFile(filename, file, 0644)

}

//todo: create_election: This API function allows the polling officer to create a new election. The function takes as input the details of the election, such as the title, start and end dates, and the list of candidates.

//todo: cast_vote: This API function allows the voter to cast their vote. The function takes as input the voter's ID and the ID of the candidate they are voting for.

//todo: end_election: This API function ends the election and stops the voting process. The function takes no input parameters.

//todo: tally_votes: This API function tallies the votes and generates the election results. The function takes no input parameters and returns the final results of the election.

//? todo: verify_voter: This API function verifies the eligibility of a voter to participate in the election. The function takes as input the voter's ID and other relevant details, and checks that they meet the eligibility criteria.

//? todo: get_voter_status: This API function retrieves the current status of a voter, such as whether they have already cast their vote or not. The function takes as input the voter's ID and returns their current status.

//? todo: cancel_vote: This API function allows a voter to cancel their vote if they made a mistake or changed their mind. The function takes as input the voter's ID and the ID of the candidate they originally voted for.

//? todo: audit_election: This API function provides an auditable trail of the entire election process, including the details of all voters, candidates, and votes. The function takes no input parameters and returns a complete audit log of the election.
