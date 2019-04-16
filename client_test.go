package main

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	addr := "127.0.0.1:9488"
	go func() {
		server := NewServer()

		srv := &http.Server{Addr: addr}
		r := mux.NewRouter()
		r.HandleFunc("/signal", server.SignalHandler)
		srv.Handler = r

		srv.ListenAndServe()
	}()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("failed: %v", err)
	}
	nodeID := discv5.PubkeyID(&key.PublicKey)

	url := url.URL{Scheme: "ws", Host: addr, Path: "/signal"}

	_, err = NewClient(nodeID, key, url)
	assert.NoError(t, err)

}

func TestClientSendReceive(t *testing.T) {
	addr := "127.0.0.1:9489"

	go func() {
		server := NewServer()

		srv := &http.Server{Addr: addr}
		r := mux.NewRouter()
		r.HandleFunc("/signal", server.SignalHandler)
		srv.Handler = r

		srv.ListenAndServe()
	}()

	url := url.URL{Scheme: "ws", Host: addr, Path: "/signal"}

	key1, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("failed: %v", err)
	}
	nodeID1 := discv5.PubkeyID(&key1.PublicKey)

	key2, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("failed: %v", err)
	}
	nodeID2 := discv5.PubkeyID(&key2.PublicKey)

	msg1 := []byte("test")
	var msg2 []byte

	c1, err := NewClient(nodeID1, key1, url)
	log.Printf("TestClientSendReceive: after c1: e: %v", err)
	assert.NoError(t, err)

	c2, err := NewClient(nodeID2, key2, url)
	log.Printf("TestClientSendReceive: after c2: e: %v", err)
	assert.NoError(t, err)

	log.Printf("TestClientSendReceive: c1 to Send c2")
	err = c1.Send(nodeID2, msg1)
	log.Printf("TestClientSendReceive: after c1 sent c2: e: %v", err)
	assert.NoError(t, err)

	msg2, err = c2.Receive()
	log.Printf("TestClientSendReceive: after c2 receive c1: msg2: %v e: %v", msg2, err)

	assert.Equal(t, msg1, msg2)
}
