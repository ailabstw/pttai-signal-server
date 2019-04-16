package main

import (
	"crypto/ecdsa"
	"log"
	"net/url"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/gorilla/websocket"
)

type Client struct {
	nodeID discv5.NodeID
	Conn   *Conn
}

/*
Send sends the normal-msg to toID.
*/
func (c *Client) Send(toID discv5.NodeID, msg []byte) error {
	signal := &signal{NodeID: toID, Msg: msg}

	err := c.Conn.WsConn.WriteJSON(signal)
	if err != nil {
		return err
	}

	return nil
}

/*
Receive receives the normal-msg from signal-server.
*/
func (c *Client) Receive() ([]byte, error) {
	signal := &signal{}
	err := c.Conn.WsConn.ReadJSON(signal)
	if err != nil {
		return nil, err
	}

	if c.nodeID != signal.NodeID {
		return nil, ErrInvalidID
	}

	return signal.Msg, nil
}

/*
NewClient init a new client and pass the challenge from the signal-server.
*/
func NewClient(nodeID discv5.NodeID, privKey *ecdsa.PrivateKey, url url.URL) (*Client, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Printf("NewClient: Error: unable to dial, e: %v", err)
		return nil, err
	}

	signal := &signal{}
	err = wsConn.ReadJSON(signal)
	if err != nil {
		log.Printf("NewClient: Error: unable to ReadJSON, e: %v\n", err)
		return nil, err
	}

	resp, err := respondChallenge(nodeID, privKey, signal)
	if err != nil {
		log.Printf("NewClient: Error: unable to respond Challenge, e: %v\n", err)
		return nil, err
	}

	err = wsConn.WriteJSON(resp)
	if err != nil {
		log.Printf("NewClient: Error: unable to WriteJSON, e: %v\n", err)
		return nil, err
	}

	err = wsConn.ReadJSON(signal)
	if err != nil {
		log.Printf("NewClient: Error: unable to ReadJSON from ack: e: %v\n", err)
		return nil, err
	}

	if signal.NodeID != nodeID {
		log.Printf("NewClient: Error: invalid id\n")
		return nil, ErrInvalidID
	}

	c := &Conn{isClosed: 0, WsConn: wsConn}

	return &Client{nodeID, c}, nil
}

func respondChallenge(nodeID discv5.NodeID, privKey *ecdsa.PrivateKey, signal *signal) (*challengeResponse, error) {

	challenge := signal.Msg
	hash := crypto.Keccak256Hash(challenge)

	sig, err := crypto.Sign(hash[:], privKey)
	if err != nil {
		return nil, err
	}

	challengeResponse := &challengeResponse{NodeID: nodeID, Signature: sig, Hash: hash}

	return challengeResponse, nil
}

func (c *Client) Close() {
	c.Conn.Close()
}
