package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discv5"
)

func TestServerIdentifyNodeID(t *testing.T) {
	server := NewServer()

	challenge := server.generateChallenge()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("failed: %v", err)
	}

	hash := crypto.Keccak256Hash(challenge)

	// sign challenge with key
	sig, err := crypto.Sign(hash[:], key)
	if err != nil {
		t.Errorf("failed: %v", err)
	}

	nodeID := discv5.PubkeyID(&key.PublicKey)
	resp := challengeResponse{
		NodeID:    nodeID,
		Hash:      hash,
		Signature: sig,
	}

	err = server.verifyNode(challenge, resp)
	if err != nil {
		t.Errorf("failed: %v", err)
	}
}
