package network

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sohelahmedjoni/pqhbevs_hac/internal/blockchain"
)

// pattern: /protocol-name/request-or-response-message/version
const blockRequest = "/block/blockreq/0.0.1"
const blockResponse = "/block/blockresp/0.0.1"

type BlockProtocol struct {
	node     *Node                               // local host
	requests map[string]*blockchain.BlockRequest // used to access request data from response handlers
	done     chan bool                           // only for demo purposes to hold main from terminating
}

func NewBlockProtocol(node *Node, done chan bool) *BlockProtocol {
	e := BlockProtocol{node: node, requests: make(map[string]*blockchain.BlockRequest), done: done}
	node.SetStreamHandler(blockRequest, e.onBlockRequest)
	node.SetStreamHandler(blockResponse, e.onBlockResponse)

	// design note: to implement fire-and-forget style messages you may just skip specifying a response callback.
	// a fire-and-forget message will just include a request and not specify a response object

	return &e
}

// remote peer requests handler
func (e *BlockProtocol) onBlockRequest(s network.Stream) {
	// get request data
	data := &blockchain.BlockRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		println("proto mark============================================================")

		log.Println(err)
		return
	}

	log.Printf("\n\n\n%s: Received block request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), spew.Sdump(data))

	valid := e.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: Sending block response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

	// send response to the request using the message string he provided
	//todo: parse block from level db
	CurrentLatestBlockMessage, err := HandleBlockReq(*data.MessageData.Liteblock)
	if err != nil {
		if err.Error() == "ALL BLOCKS ARE UP TO DATE" {
			log.Println("ALL BLOCKS ARE UP TO DATE")
			data.Message = "YOU ARE UP TO DATE"
		} else {
			panic(err)
		}
	} else {

	}
	resp := &blockchain.BlockResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false, CurrentLatestBlockMessage),
		Message:     data.Message}

	// sign the data
	signature, err := e.node.signProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature
	//? handle error
	ctx := context.Background()
	ok := e.node.sendProtoMessage(ctx, s.Conn().RemotePeer(), blockResponse, resp)

	if ok {
		log.Printf("%s: Block response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
	e.done <- true
}

// remote block response handler
func (e *BlockProtocol) onBlockResponse(s network.Stream) {

	data := &blockchain.BlockResponse{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)

		log.Println(err)
		return
	}

	// authenticate message content
	valid := e.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate request data and remove it if found
	req, ok := e.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(e.requests, data.MessageData.Id)
	} else {
		log.Println("Failed to locate request data project for response")
		return
	}
	if data.Message == "YOU ARE UP TO DATE" {
		log.Println("!!YOU ARE UP TO DATE")
		e.done <- true
		return
	}
	if req.Message != data.Message {

		log.Fatalln("Expected block to respond with request message ", req.Message, " but got: ", data.Message)
	}

	log.Printf("\n\n\n%s: ~~Received block response from %s. Message id:%s. \nMessage: %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id, spew.Sdump(data.MessageData))
	//todo: check and save block to blockchain level db
	if req.MessageData.Block != nil {
		if req.MessageData.Block.Header.BlockID > blocky.Block.Header.BlockID {
			log.Printf("Block is valid. Adding to blockchain")
			if !blockchain.AddBlockAsLatest(req.MessageData.Block, req.MessageData.Binarydata) {
				log.Printf("Block is invalid. Discarding")
			}
		}

	} else {
		log.Printf("Block is invalid. Discarding")
	}
	e.done <- true
}

// SendBlockData sends latest lite block to peer
// so they can verify we are synced or he can request the block from us
// index is to set the record needed to be sent
// returns true if the message was sent successfully
func (e *BlockProtocol) SendBlockData(ctx context.Context, peer peer.AddrInfo) bool {
	log.Printf("%s: Sending latest lite block to: %s....", e.node.ID(), peer.ID)
	//todo: parse block from level db
	CurrentLatestBlock := blocky.getLiteBlock()
	msg := blockchain.MessageData{
		Liteblock: &CurrentLatestBlock,
	}
	// create message data
	req := &blockchain.BlockRequest{
		MessageData: e.node.NewMessageData(uuid.New().String(), false, &msg),
		Message:     fmt.Sprintf("Block from %s", e.node.ID())}

	signature, err := e.node.signProtoMessage(req)
	if err != nil {
		log.Println("failed to sign message")
		return false
	}

	// add the signature to the message
	req.MessageData.Sign = signature

	ok := e.node.sendProtoMessage(ctx, peer.ID, blockRequest, req)

	if !ok {
		return false
	}

	// store request so response handler has access to it
	e.requests[req.MessageData.Id] = req
	log.Printf("%s: ~~Block to: %s was sent. Message Id: %s, \nMessage: %s", e.node.ID(), peer.ID, req.MessageData.Id, spew.Sdump(req.MessageData))
	return true
}

/*
func (b *Block) getDigest() string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprint(b.Header.BlockID)))
	hash.Write([]byte(b.Header.PreHash))
	hash.Write([]byte(b.Header.PostHash))
	hash.Write([]byte(fmt.Sprint(b.Header.BlockID)))
	hash.Write([]byte(fmt.Sprint(b.Header.ElectionID)))
	hash.Write([]byte(b.Header.Pscode))
	hash.Write([]byte(fmt.Sprint(b.Header.Timestamp)))

	hash.Write([]byte(b.CandidateStruct.Digest))
	hash.Write([]byte(b.ContractStruct.getDigest()))
	hash.Write([]byte(b.ElectionCommission.Digest))
	hash.Write([]byte(b.ElectionTime.GenerateDigest()))
	hash.Write([]byte(b.SpendedToken.Digest))
	for i := 0; i < len(b.PoBannedList.PoHash[i]); i++ {
		hash.Write([]byte(b.PoBannedList.PoHash[i]))
	}
	for i := 0; i < len(b.PollingOfficer); i++ {
		hash.Write([]byte(b.PollingOfficer[i].Digest))
	}
	for i := 0; i < len(b.ReturningOfficer); i++ {
		hash.Write([]byte(b.ReturningOfficer[i].Digest))
	}
	return base58.Encode(hash.Sum(nil))
}
*/
