package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/mr-tron/base58"
	"github.com/sohelahmedjoni/pqhbevs_hac/pkg/utils"
)

//? todo: verify_voter: This API function verifies the eligibility of a voter to participate in the election. The function takes as input the voter's ID and other relevant details, and checks that they meet the eligibility criteria.

func verify_voter() {

}

// ? todo: get_voter_status: This API function retrieves the current status of a voter, such as whether they have already cast their vote or not. The function takes as input the voter's ID and returns their current status.
func get_voter_status() {

}

// will be interfaced to VoterStruct struct
func GenerateVoterList(config, pscode, seed string) VoterStruct {
	var voterList VoterStruct

	voterList.Pscode = pscode
	voterList.Seed = seed
	//todo add custom port
	if config != "1" {
		bit, err := utils.Get_VoterList(pscode, seed)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bit, &voterList)
		if err != nil {
			panic(err)
		}
		fmt.Println("voterList: ")
		fmt.Println(spew.Sdump(voterList))

	} else {
		//parse info from configuration
		configFile, err := os.Open("config/voterlist.json")
		if err != nil {
			panic(err)
		}
		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&voterList); err != nil {
			panic(err)
		}
	}

	digest := sha256.New()
	for _, v := range voterList.Voters {
		digest.Write([]byte(fmt.Sprintf("%v", v)))
	}
	digest.Write([]byte(fmt.Sprintf("%v", voterList.Pscode)))
	digest.Write([]byte(fmt.Sprintf("%v", voterList.Seed)))

	voterList.Digest = base58.Encode(digest.Sum(nil))

	return voterList
}
