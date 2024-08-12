package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

// all element of block are signed by election commission & first returning officer
// todo: make it modular so returning officer can be changed
func (b *Block) CreateGenesisBlock(startTime, duration, config, pscode, seed string, db *leveldb.DB) {
	// declaring all the required instance of blockcs

	var vEC ElectionCommission
	var vcandidate CandidateStruct
	var vReturningOfficer []ReturningOfficer
	var vPollingOfficer []PollingOfficer
	var voterList VoterStruct
	var vContract ContractStruct
	var vHeader BlockHeader
	var vSpendedToken SpendedToken
	var PoBannedListStruct PoBannedList
	var ElectionTimeStruct = ElectionTime{Start: startTime, Duration: duration}
	ElectionTimeStruct.Digest = ElectionTimeStruct.GenerateDigest()

	//* Generating election commission
	vEC.GenerateElection(db)
	b.ElectionCommission = &vEC

	//* Generating returning officer
	vReturningOfficer = GenerateReturningOfficer("EC_private_key", db)
	var blob []byte
	for i := 0; i < len(vReturningOfficer); i++ {
		blob = utils.GetAuthorize(vReturningOfficer[i].Digest, vEC.PrivateKey, db)
		vReturningOfficer[i].Signature = utils.BlobhashGeneration(blob)
		err := db.Put([]byte(vReturningOfficer[i].Signature), blob, nil)
		if err != nil {
			log.Panic(err)
		}
		err = db.Put([]byte("signature_"+vReturningOfficer[i].Id), blob, nil)
		if err != nil {
			log.Panic(err)
		}

		blob, err = vReturningOfficer[i].XXX_Marshal(nil, true)
		check(err)
		err = db.Put([]byte(fmt.Sprint(vReturningOfficer[i].Digest)), blob, nil)
		check(err)
		err = db.Put([]byte(fmt.Sprint(vReturningOfficer[i].Id)), blob, nil)
		check(err)
	}

	//* Generating polling officer
	vPollingOfficer = GeneratePollingOfficer(vReturningOfficer[0].Id, db)

	for i := 0; i < len(vPollingOfficer); i++ {
		//? fix returning officer index
		blob = utils.GetAuthorize(vPollingOfficer[i].Digest, vReturningOfficer[0].Id, db)
		vPollingOfficer[i].Signature = utils.BlobhashGeneration(blob)
		err := db.Put([]byte(vPollingOfficer[i].Signature), blob, nil)
		if err != nil {
			log.Panic(err)
		}
		err = db.Put([]byte("signature_"+vPollingOfficer[i].Id), blob, nil)
		if err != nil {
			log.Panic(err)
		}
		blob, err = vPollingOfficer[i].XXX_Marshal(nil, true)
		check(err)
		err = db.Put([]byte(fmt.Sprint(vPollingOfficer[i].Digest)), blob, nil)
		check(err)
		err = db.Put([]byte(fmt.Sprint(vPollingOfficer[i].Id)), blob, nil)
		check(err)
	}
	// election time authorization
	blob = utils.GetAuthorize(ElectionTimeStruct.Digest, vReturningOfficer[0].Id, db)
	ElectionTimeStruct.ETproof = utils.BlobhashGeneration(blob)
	err := db.Put([]byte(ElectionTimeStruct.ETproof), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("election_time_signature"), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	blob, err = ElectionTimeStruct.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(ElectionTimeStruct.Digest), blob, nil)
	check(err)
	err = db.Put([]byte("election_time"), blob, nil)
	check(err)

	//! tree mod
	//* assigning intial spended token
	tokenSet := treeset.NewWithStringComparator()
	tokenSet.Add("0")
	TokenBytes, err := tokenSet.ToJSON()
	check(err)

	// token save in db
	vSpendedToken.Token = utils.BlobhashGeneration(TokenBytes)
	err = db.Put([]byte(vSpendedToken.Token), TokenBytes, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("spended_token_tree"), TokenBytes, nil)
	if err != nil {
		log.Panic(err)
	}
	// generating digest of token

	// Generating proof of spended token
	//? fix returning officer index
	blob, err = vSpendedToken.XXX_Marshal(nil, true)
	check(err)

	vSpendedToken.Digest = utils.BlobhashGeneration(blob)
	err = db.Put([]byte(vSpendedToken.Digest), blob, nil)
	check(err)
	err = db.Put([]byte("spended_token"), blob, nil)
	check(err)

	// ban list
	PoBannedListStruct = PoBannedList{}
	PoBannedListStruct.PoHash = append(PoBannedListStruct.PoHash, "0")
	PoBannedListStruct.PoPoa = make([]string, len(PoBannedListStruct.PoHash))

	// Generating contract
	vcandidate = GenerateCandidate()

	//? fix returning officer index
	blob = utils.GetAuthorize(vcandidate.Digest, vReturningOfficer[0].Keys.PrivateKey, db)
	vcandidate.Signature = utils.BlobhashGeneration(blob)
	err = db.Put([]byte(vcandidate.Signature), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("candidate_signature"), blob, nil)
	if err != nil {
		log.Panic(err)
	}

	voterList = GenerateVoterList(config, pscode, seed)
	//? fix returning officer index
	blob = utils.GetAuthorize(voterList.Digest, vReturningOfficer[0].Keys.PrivateKey, db)

	voterList.Signature = utils.BlobhashGeneration(blob)
	err = db.Put([]byte(voterList.Signature), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("voterlist_Signature"), blob, nil)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < len(vReturningOfficer); i++ {
		vContract.Command = append(vContract.Command,
			&Command{AuthorizeBy: "EC_private_key", Action: "ASSIGN", Entity: vPollingOfficer[i].Id})
		vContract.Digest = append(vContract.Digest, &Digest{Digest: vContract.Command[i].GenerateDigest()})
	}

	for i := 0; i < len(vPollingOfficer); i++ {
		vContract.Command = append(vContract.Command,
			&Command{AuthorizeBy: vReturningOfficer[0].Id, Action: "ASSIGN", Entity: vPollingOfficer[i].Id})
		vContract.Digest = append(vContract.Digest, &Digest{Digest: vContract.Command[i].GenerateDigest()})
	}

	//need contract stuff
	// for i := 0; i < len(vContract.Command); i++ {
	// 	vContract.Command = append(vContract.Command, &Command{AuthorizeBy: vReturningOfficer[i].Id, Action: "ASSIGN", Entity: "candidate"})
	// }

	// vContract.ProofOfCommand = make([]*ProofOfCommand, len(vContract.Digest))
	for i := 0; i < len(vContract.Digest); i++ {
		blob = utils.GetAuthorize(vContract.Digest[i].Digest, vContract.Command[i].AuthorizeBy, db)
		vContract.ProofOfCommand = append(vContract.ProofOfCommand, &ProofOfCommand{ProofOfCommand: utils.BlobhashGeneration(blob)})
		err = db.Put([]byte(vContract.ProofOfCommand[i].ProofOfCommand), blob, nil)
		if err != nil {
			log.Panic(err)
		}
		err = db.Put([]byte("proof_of_command_"+fmt.Sprint(i+1)), blob, nil)
		if err != nil {
			log.Panic(err)
		}

	}
	blob, err = vContract.XXX_Marshal(nil, true)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("contract"), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte(vContract.getDigest()), blob, nil)
	if err != nil {
		log.Panic(err)
	}

	blob, err = ElectionTimeStruct.XXX_Marshal(nil, true)
	check(err)

	err = db.Put([]byte("election_time"), blob, nil)
	check(err)
	err = db.Put([]byte(fmt.Sprint(ElectionTimeStruct.Digest)), blob, nil)
	check(err)

	blob, err = vEC.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(fmt.Sprint(vEC.Digest)), blob, nil)
	check(err)
	err = db.Put([]byte("election_commission"), blob, nil)
	check(err)

	blob, err = vcandidate.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(fmt.Sprint(vcandidate.Digest)), blob, nil)
	check(err)
	err = db.Put([]byte("candidate_struct"), blob, nil)
	check(err)

	blob, err = voterList.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(fmt.Sprint(voterList.Digest)), blob, nil)
	check(err)
	err = db.Put([]byte("voter_list"), blob, nil)
	check(err)

	for i := 0; i < len(PoBannedListStruct.GetPoHash()); i++ {
		blob = utils.GetAuthorize(PoBannedListStruct.PoHash[i], vReturningOfficer[0].GetId(), db)
		PoBannedListStruct.PoHash[i] = utils.BlobhashGeneration(blob)
		err := db.Put([]byte(PoBannedListStruct.PoHash[i]), blob, nil)
		if err != nil {
			log.Panic(err)
		}
	}
	blob, err = PoBannedListStruct.XXX_Marshal(nil, true)
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("PoBannedList"), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	//genrating light block
	var lb LiteBlock
	for i := 0; i < len(vReturningOfficer); i++ {
		lb.ReturningOfficer = append(lb.ReturningOfficer, vReturningOfficer[i].Digest)
	}
	for i := 0; i < len(vPollingOfficer); i++ {
		lb.PollingOfficer = append(lb.PollingOfficer, vPollingOfficer[i].Digest)
	}
	lb.ElectionCommission = vEC.Digest
	lb.ElectionTime = ElectionTimeStruct.Digest
	lb.CandidateStruct = vcandidate.Digest
	lb.VoterStruct = voterList.Digest
	lb.SpendedToken = vSpendedToken.Digest
	lb.PoBannedList = "PoBannedList"
	lb.ContractStruct = vContract.getDigest()

	vHeader = GenerateHeaderlb(&lb)

	//? fix returning officer index
	blob = utils.GetAuthorize(vHeader.PostHash, vReturningOfficer[0].Id, db)
	vHeader.ProofOfBlock = utils.BlobhashGeneration(blob)
	err = db.Put([]byte(vHeader.ProofOfBlock), blob, nil)
	if err != nil {
		log.Panic(err)
	}
	lb.Header = &vHeader
	blob, err = lb.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte("genesis_lite_block"), blob, nil)
	check(err)
	err = db.Put([]byte("lite_block_0"), blob, nil)
	check(err)
	err = db.Put([]byte("latest_lite_block"), blob, nil)
	check(err)

	b.ReturningOfficer = make([]*ReturningOfficer, len(vReturningOfficer))
	//generating heavy block
	for i := 0; i < len(vReturningOfficer); i++ {
		b.ReturningOfficer[i] = &vReturningOfficer[i]
	}
	b.PollingOfficer = make([]*PollingOfficer, len(vPollingOfficer))
	for i := 0; i < len(vPollingOfficer); i++ {
		b.PollingOfficer[i] = &vPollingOfficer[i]
	}
	b.VoterStruct = &voterList
	b.ContractStruct = &vContract
	b.ElectionTime = &ElectionTimeStruct
	b.Header = &vHeader
	b.SpendedToken = &vSpendedToken
	// spew.Dump(ElectionTimeStruct)
	//generating jsonfile of the genesis block
	file, _ := json.MarshalIndent(b, "", " ")
	filename := "./BLOCK/genesis.json"
	os.WriteFile(filename, file, fs.ModePerm)
	// spew.Dump(lb)
	// saving the genesis in db *level.DB
	blob, err = b.XXX_Marshal(nil, true)
	check(err)
	err = db.Put([]byte(fmt.Sprintf("block_%v", b.Header.BlockID)), blob, nil)
	check(err)
	err = db.Put([]byte("genesis"), blob, nil)
	check(err)
	err = db.Put([]byte("genesis_block"), blob, nil)
	check(err)
	err = db.Put([]byte("latest"), blob, nil)
	check(err)

}

// This function parse json from
func (b *Block) Parsefromjson(file string) {
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("opening json file: %s\n", err.Error())
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&b); err != nil {
		if err != nil {
			panic(err)
		}
	}
}

// * This function will check voter va.Idity
func (lb *LiteBlock) Checkblock(key []byte) error {
	if Generate_block_hash(*lb) != lb.Header.PostHash {
		fmt.Printf("")
	}
	return nil
	//TODO: Verify the signature of the Election Commission
	//TODO: Verify the signature of the Returning Officers
	//TODO: Verify the signature of the Polling Officers
	//TODO: Verify the signature of the Voter List
	//TODO: Verify the signature of the Contract
	//TODO: Check for any banned Polling Officers
	//TODO: Check for any banned Polling Officer's Proof of Authority
	//TODO: Check for any other necessary va.Idations based on the requirements of the voting system
}

// * This function will return hevy block
func (lb *LiteBlock) GetFullBlock() Block {
	// TODO: Get the full block from the database
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

	var b Block

	var vEC ElectionCommission
	var vcandidate CandidateStruct
	var vReturningOfficer []ReturningOfficer
	var vPollingOfficer []PollingOfficer
	var voterList VoterStruct
	var vContract ContractStruct
	var vSpendedToken SpendedToken
	var PoBannedListStruct PoBannedList
	var ElectionTimeStruct ElectionTime
	// println("liteblock:")
	//! spew.Dump(lb)

	blob, err := db.Get([]byte(lb.ElectionCommission), nil)
	if err != nil {
		log.Panic(err)
	}
	vEC.XXX_Unmarshal(blob)
	blob, err = db.Get([]byte(lb.CandidateStruct), nil)
	if err != nil {
		log.Panic(err)
	}
	vcandidate.XXX_Unmarshal(blob)

	for i := 0; i < len(lb.ReturningOfficer); i++ {
		blob, err = db.Get([]byte(lb.ReturningOfficer[i]), nil)
		if err != nil {
			log.Panic(err)
		}
		var vReturningOfficerTemp ReturningOfficer
		vReturningOfficerTemp.XXX_Unmarshal(blob)
		vReturningOfficer = append(vReturningOfficer, vReturningOfficerTemp)
	}

	for i := 0; i < len(lb.PollingOfficer); i++ {
		blob, err = db.Get([]byte(lb.PollingOfficer[i]), nil)
		if err != nil {
			log.Panic(err)
		}
		var vPollingOfficerTemp PollingOfficer
		vPollingOfficerTemp.XXX_Unmarshal(blob)
		vPollingOfficer = append(vPollingOfficer, vPollingOfficerTemp)
	}

	blob, err = db.Get([]byte(lb.VoterStruct), nil)
	if err != nil {
		log.Panic(err)
	}
	voterList.XXX_Unmarshal(blob)

	blob, err = db.Get([]byte(lb.ContractStruct), nil)
	if err != nil {
		log.Panic(err)
	}
	vContract.XXX_Unmarshal(blob)

	blob, err = db.Get([]byte("spended_token"), nil)
	if err != nil {
		log.Panic(err)
	}
	vSpendedToken.XXX_Unmarshal(blob)

	b.Header = lb.Header
	b.ElectionTime = &ElectionTimeStruct
	b.CandidateStruct = &vcandidate
	b.VoterStruct = &voterList
	b.ElectionCommission = &vEC
	b.ReturningOfficer = make([]*ReturningOfficer, len(lb.ReturningOfficer))
	b.PollingOfficer = make([]*PollingOfficer, len(lb.PollingOfficer))
	b.ContractStruct = &vContract
	b.PoBannedList = &PoBannedListStruct
	b.SpendedToken = &vSpendedToken
	b.Token = lb.Token
	b.Vote = lb.Vote

	for i := 0; i < len(lb.ReturningOfficer); i++ {
		b.ReturningOfficer[i] = &vReturningOfficer[i]
	}
	for i := 0; i < len(lb.PollingOfficer); i++ {
		b.PollingOfficer[i] = &vPollingOfficer[i]
	}
	// fmt.Printf("Full Block: \n")
	// jbyte, err := json.MarshalIndent(b, "", "    ")
	// fmt.Printf("%+v", string(jbyte))
	return b
}

/*

	lb.ElectionCommission = vEC.Digest
	lb.ElectionTime = ElectionTimeStruct.Digest
	lb.CandidateStruct = vcandidate.Digest
	lb.VoterStruct = voterList.Digest
	lb.SpendedToken = vSpendedToken.Digest
	lb.PoBannedList = "PoBannedList"
	lb.ContractStruct = vContract.getDigest()


	Header               *BlockHeader `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
    ElectionTime         string       `protobuf:"bytes,2,opt,name=ElectionTime,proto3" json:"ElectionTime,omitempty"`
    CandidateStruct      string       `protobuf:"bytes,3,opt,name=CandidateStruct,proto3" json:"CandidateStruct,omitempty"`
    VoterStruct          string       `protobuf:"bytes,4,opt,name=VoterStruct,proto3" json:"VoterStruct,omitempty"`
    ElectionCommission   string       `protobuf:"bytes,5,opt,name=ElectionCommission,proto3" json:"ElectionCommission,omitempty"`
    ContractStruct       string       `protobuf:"bytes,8,opt,name=ContractStruct,proto3" json:"ContractStruct,omitempty"`
    PoBannedList         string       `protobuf:"bytes,9,opt,name=PoBannedList,proto3" json:"PoBannedList,omitempty"`
    SpendedToken         string       `protobuf:"bytes,10,opt,name=spendedToken,proto3" json:"spendedToken,omitempty"`
    Token                string       `protobuf:"bytes,11,opt,name=Token,proto3" json:"Token,omitempty"`
    Vote
    PollingOfficer       []string     `protobuf:"bytes,6,rep,name=PollingOfficer,proto3" json:"PollingOfficer,omitempty"`
    ReturningOfficer     []string     `protobuf:"bytes,7,rep,name=ReturningOfficer,proto3" json:"ReturningOfficer,omitempty"`
*/

// It'll generate hash of the election time
func (v *ElectionTime) GenerateDigest() string {

	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", v.Start)))
	digest.Write([]byte(fmt.Sprintf("%v", v.Duration)))
	return base58.Encode(digest.Sum(nil))
}

// It'll generate hash of the CURRENT BLOCK
// It'll take hash from each unit of the block
// and generate a unique digest of the block
func Generate_block_hash(lb LiteBlock) string {

	digest := sha256.New()
	digest.Write([]byte(fmt.Sprintf("%v", lb.Header.BlockID)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.Header.PreHash)))
	digest.Write([]byte(fmt.Sprintf("%v", lb.Header.Pscode)))

	//! timestamp caution and need to test throughly
	digest.Write([]byte(fmt.Sprintf("%v", lb.Header.Timestamp)))

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

	//! need to add PObanned list
	//! need to add Vote
	return base58.Encode(digest.Sum(nil))
}

func GetBlock(BlockID uint64) []byte {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	parent = path.Join(wd, "BLOCK")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Working Directory: ", parent)

	//terminating if the directory does not exist
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		log.Printf("Directory does not exist, Make sure your working directory is correct\n")
		panic(err)
	}

	db, err := leveldb.OpenFile(path.Join(parent, "BLOCK"), nil)
	if err != nil {
		log.Printf("Failed to open database, Make sure your working directory is correct\n")
		panic(err)
	}
	defer db.Close()
	block, err := db.Get([]byte(fmt.Sprintf("block_%v", BlockID)), nil)
	if err != nil {
		panic(err)
	}

	return block
}
func GetLatestBlock() []byte {
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
	block, err := db.Get([]byte("latest_lite_block"), nil)
	if err != nil {
		panic(err)
	}

	return block
}

func GetSelectedRecordFromBlock(fields []bool, blockID string) *MessageData {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := path.Join(wd, "BLOCK")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Working Directory: ", parent)

	db, err := leveldb.OpenFile(parent, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	blob, err := db.Get([]byte(blockID), nil)
	check(err)
	block := new(LiteBlock)
	block.XXX_Unmarshal(blob)
	// spew.Dump(block)
	fullblock := new(Block)
	fullblock.Header = block.Header
	binProofOfBlock, err := db.Get([]byte(block.Header.ProofOfBlock), nil)
	check(err)

	bindata := new(BinaryData)
	// ElectionTime get
	if fields[1] {
		blob, err = db.Get([]byte(block.ElectionTime), nil)
		check(err)
		fullblock.ElectionTime = new(ElectionTime)
		fullblock.ElectionTime.XXX_Unmarshal(blob)
		bindata.ETproof, err = db.Get([]byte(fullblock.ElectionTime.ETproof), nil)
		check(err)

	}
	if fields[2] {
		// CANDIDATE get
		blob, err = db.Get([]byte(block.CandidateStruct), nil)
		check(err)
		fullblock.CandidateStruct = new(CandidateStruct)
		fullblock.CandidateStruct.XXX_Unmarshal(blob)
		bindata.CandidateSignature, err = db.Get([]byte(fullblock.CandidateStruct.Signature), nil)
		check(err)
	}
	if fields[3] {
		// VOTER get
		blob, err = db.Get([]byte(block.VoterStruct), nil)
		check(err)
		fullblock.VoterStruct = new(VoterStruct)
		fullblock.VoterStruct.XXX_Unmarshal(blob)
		bindata.VotersSignature, err = db.Get([]byte(fullblock.VoterStruct.Signature), nil)
		check(err)
	}
	if fields[4] {
		// GET ELECTION COMMISSION
		blob, err = db.Get([]byte(block.ElectionCommission), nil)
		check(err)
		fullblock.ElectionCommission = new(ElectionCommission)
		fullblock.ElectionCommission.XXX_Unmarshal(blob)
	}
	if fields[5] {
		// Polling officer get
		fullblock.PollingOfficer = make([]*PollingOfficer, len(block.PollingOfficer))
		for i := 0; i < len(block.PollingOfficer); i++ {
			blob, err = db.Get([]byte(block.PollingOfficer[i]), nil)
			check(err)
			fullblock.PollingOfficer[i] = new(PollingOfficer)
			fullblock.PollingOfficer[i].XXX_Unmarshal(blob)
			blob, err = db.Get([]byte(fullblock.PollingOfficer[i].Signature), nil)
			check(err)
			bindata.PollingOfficerSignature = append(bindata.PollingOfficerSignature, blob)
		}
	}
	if fields[6] {
		// Returning officer get
		fullblock.ReturningOfficer = make([]*ReturningOfficer, len(block.ReturningOfficer))
		for i := 0; i < len(block.ReturningOfficer); i++ {
			blob, err = db.Get([]byte(block.ReturningOfficer[i]), nil)
			check(err)
			fullblock.ReturningOfficer[i] = new(ReturningOfficer)
			fullblock.ReturningOfficer[i].XXX_Unmarshal(blob)
			blob, err = db.Get([]byte(fullblock.ReturningOfficer[i].Signature), nil)
			check(err)
			bindata.ReturningOfficerSignature = append(bindata.ReturningOfficerSignature, blob)

		}
	}
	if fields[7] {
		// Contract get
		blob, err = db.Get([]byte(block.ContractStruct), nil)
		check(err)
		fullblock.ContractStruct = new(ContractStruct)
		fullblock.ContractStruct.XXX_Unmarshal(blob)
		for i := 0; i < len(fullblock.ContractStruct.Digest); i++ {
			blob, err = db.Get([]byte(fullblock.ContractStruct.ProofOfCommand[i].ProofOfCommand), nil)
			check(err)
			bindata.ContractProofOfCommand = append(bindata.ContractProofOfCommand, blob)
		}
		check(err)

	}
	if fields[8] {
		// get po banned list
		blob, err = db.Get([]byte(block.PoBannedList), nil)
		check(err)
		fullblock.PoBannedList = new(PoBannedList)
		for i := 0; i < len(fullblock.PoBannedList.GetPoHash()); i++ {
			blob, err := db.Get([]byte(fullblock.PoBannedList.PoHash[i]), nil)
			if err != nil {
				log.Panic(err)
			}
			bindata.PoBannedList = append(bindata.PoBannedList, blob)

		}

	}
	if fields[9] {
		// get Spended token
		blob, err = db.Get([]byte(block.SpendedToken), nil)
		check(err)
		fullblock.SpendedToken = new(SpendedToken)
		fullblock.SpendedToken.XXX_Unmarshal(blob)
	}

	// get token
	fullblock.Token = block.Token
	// get vote
	fullblock.Vote = block.Vote

	return &MessageData{Block: fullblock, ProofOfBlock: binProofOfBlock, Binarydata: bindata}
}

/*
BlockHeader Header = 0;
ElectionTime ElectionTime = 1;
CandidateStruct CandidateStruct = 2;
VoterStruct VoterStruct = 3;
ElectionCommission ElectionCommission = 4;
repeated  PollingOfficer PollingOfficer = 5;
repeated  ReturningOfficer ReturningOfficer =  6;
ContractStruct ContractStruct =   7;
PoBannedList PoBannedList =  8;
spendedToken spendedToken  =  9;

			string Token =  10;
			string Vote =  11;
*/
