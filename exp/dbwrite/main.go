package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
)

const (
	count  = 30
	dbpath = "db"
)

func main() {

	db, err := pebble.Open(dbpath, nil)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	store := store.New(db, JSONSerializer{})

	peers := generatePeers(count)

	for _, p := range peers {

		data, _ := json.Marshal(p)
		fmt.Printf("%s\n", data)

		err = store.SavePeer(p)
		if err != nil {
			log.Fatalf("could not save peer: %s", err)
		}
	}

	functions := generateFunctions(count)
	for _, f := range functions {

		data, _ := json.Marshal(f)
		fmt.Printf("%s\n", data)

		err = store.SaveFunction(f.CID, f)
		if err != nil {
			log.Fatalf("could not save function: %s", err)
		}
	}
}

func generatePeers(count int) []blockless.Peer {

	peers := make([]blockless.Peer, count)
	for i := 0; i < count; i++ {

		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 0, rand.New(rand.NewSource(int64(i))))
		if err != nil {
			log.Fatalf("could not generate key pair: %s", err)
		}

		id, err := peer.IDFromPrivateKey(priv)
		if err != nil {
			log.Fatalf("could not generate peer ID: %s", err)
		}

		addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/9010/p2p/" + id.String())
		if err != nil {
			log.Fatalf("could not parse multiaddress: %s", err)
		}

		p := blockless.Peer{
			ID:        id,
			MultiAddr: addr.String(),
			AddrInfo: peer.AddrInfo{
				ID:    id,
				Addrs: []multiaddr.Multiaddr{addr},
			},
		}

		peers[i] = p
	}

	return peers
}

func generateFunctions(count int) []blockless.FunctionRecord {

	functions := make([]blockless.FunctionRecord, count)
	for i := 0; i < count; i++ {

		fn := blockless.FunctionRecord{
			CID: fmt.Sprintf("dummy-cid-%v", i),
			URL: fmt.Sprintf("https://example.com/dummy-url-%v", i),
		}

		functions[i] = fn
	}

	return functions
}
