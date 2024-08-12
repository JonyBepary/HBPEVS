package network

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/sohelahmedjoni/pqhbevs_hac/internal/blockchain"

	ma "github.com/multiformats/go-multiaddr"
)

type Config struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
	ApiPort          int
}

type block struct {
	blockchain.Block
	blockchain.LiteBlock
}

func (b *block) getBlock() blockchain.Block {
	return b.Block
}
func (b *block) getLiteBlock() blockchain.LiteBlock {
	return b.LiteBlock
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

func HandleBlockReq(blk blockchain.LiteBlock) (*blockchain.MessageData, error) {

	fields := make([]bool, 0)
	if blk.Header.GetBlockID() < blocky.LiteBlock.Header.GetBlockID() {
		fields = append(fields, true)
		if blk.ElectionTime != blocky.LiteBlock.ElectionTime {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		if blk.CandidateStruct != blocky.LiteBlock.CandidateStruct {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		if blk.VoterStruct != blocky.LiteBlock.VoterStruct {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		if blk.ElectionCommission != blocky.LiteBlock.ElectionCommission {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		//todo: more secure hashbased checking method needed
		if len(blk.PollingOfficer) > len(blocky.LiteBlock.PollingOfficer) {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		//todo: more secure hashbased checking method needed
		if len(blk.ReturningOfficer) > len(blocky.LiteBlock.ReturningOfficer) {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}

		if blk.ContractStruct != blocky.LiteBlock.ContractStruct {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		if blk.PoBannedList != blocky.LiteBlock.PoBannedList {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		if blk.SpendedToken != blocky.LiteBlock.SpendedToken {
			fields = append(fields, true)
		} else {
			fields = append(fields, false)
		}
		log.Printf("fields length: %d\n", len(fields))
		b := blockchain.GetSelectedRecordFromBlock(fields, fmt.Sprintf("block_%v", blk.Header.GetBlockID()+1))

		return b, nil
	} else {
		b := &blockchain.MessageData{
			Block:        &blockchain.Block{},
			Liteblock:    &blockchain.LiteBlock{},
			ProofOfBlock: []byte{},
			Binarydata:   &blockchain.BinaryData{},
		}
		return b, fmt.Errorf("ALL BLOCKS ARE UP TO DATE")
	}
}

var blocky block

func NodeInstance(cfg *Config) {

	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)
	bit := blockchain.GetLatestBlock()
	err := blocky.LiteBlock.XXX_Unmarshal(bit)
	if err != nil {
		fmt.Println(err)
	}
	//! spew.Dump(blocky.LiteBlock)
	blocky.Block = blocky.LiteBlock.GetFullBlock()

	ctx := context.Background()

	done := make(chan bool, 1)

	// Make a random hosts
	host := makeRandomNode(cfg.ListenPort, done)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())

	peerChan := initMDNS(host, cfg.RendezvousString)
	for { // allows multiple peers to join
		peer := <-peerChan // will block untill we discover a peer

		fmt.Println("Found peer:", peer, ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		if host.SendBlockData(ctx, peer) {
			fmt.Println("Connected to:", peer)
			host.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)
		} else {
			fmt.Println("Connection failed")
		}

	}

	// block until all responses have been processed
	// for i := 0; i < 4; i++ {
	// 	<-done
	// }
}

// helper method - create a lib-p2p host to listen on a port
func makeRandomNode(port int, done chan bool) *Node {
	// Ignoring most errors for brevity
	// See echo example for more details and better implementation
	priv, _, _ := crypto.GenerateKeyPair(crypto.Ed25519, 256)
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	return NewNode(host, done)
}
