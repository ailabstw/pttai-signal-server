package signalserver

import (
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
		t.Errorf("failed generate key1: %v", err)
	}
	nodeID1 := discv5.PubkeyID(&key1.PublicKey)

	key2, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("failed generate key2: %v", err)
	}
	nodeID2 := discv5.PubkeyID(&key2.PublicKey)

	c1, err := NewClient(nodeID1, key1, url)
	t.Logf("TestClientSendReceive: after c1: e: %v", err)
	assert.NoError(t, err)

	c2, err := NewClient(nodeID2, key2, url)
	t.Logf("TestClientSendReceive: after c2: e: %v", err)
	assert.NoError(t, err)

	t.Logf("TestClientSendReceive: c1 to Send c2")
	msg1 := []byte("test")
	err = c1.Send(nodeID2, msg1, nil)
	t.Logf("TestClientSendReceive: after c1 sent c2: e: %v", err)
	assert.NoError(t, err)

	sig2, err := c2.Receive()
	t.Logf("TestClientSendReceive: after c2 receive c1: sig2: %v e: %v", sig2, err)

	assert.Equal(t, msg1, sig2.Msg)

}
