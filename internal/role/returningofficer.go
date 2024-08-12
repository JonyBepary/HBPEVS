package role

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohelahmedjoni/pqhbevs_hac/internal/blockchain"
)

func get_polling_officer_sign(c *gin.Context) {
	next_block := new(blockchain.Block)
	next_block.Parsefromjson("config/config.json")

	// next_block.PollingOfficer.PoPublicKey = base64_blob_to_encoded_string(PollingOfficer_publicKey)
	// next_block.PollingOfficer.Digest = Po_digestGeneration(next_block)
	// // fmt.Println("!@#!@#")
	// generate_signature_of_a_string_to_file(next_block.PollingOfficer.Digest, ReturningOfficer_privateKey, PollingOfficer_signature)

	c.String(http.StatusAccepted, fmt.Sprintf("done!!\n"))
}

func get_returning_officer_sign(c *gin.Context) {
	// next_block := new(LiteBlock)
	// if !parsefromjson("config/config.json", next_block) {
	// 	fmt.Println("config/config.json")
	// 	fmt.Println("config.json not found!!!")
	// 	return
	// }

	// next_block.ReturningOfficer.RoPublicKey = base64_blob_to_encoded_string(ReturningOfficer_publicKey)
	// next_block.ReturningOfficer.Digest = Ro_digestGeneration(next_block)
	// fmt.Println("!@#!@#")
	// generate_signature_of_a_string_to_file(next_block.ReturningOfficer.Digest, EC_privateKey, ReturningOfficer_signature)
	c.String(http.StatusAccepted, fmt.Sprintf("Done!!!\n"))
}
