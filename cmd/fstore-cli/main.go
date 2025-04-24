package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/store/codec"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	const (
		counterCID    = "bafybeianyomxw63ikxrrlhpsew6sgje4q6pud3xibxqiwghfeprcrqilxe"
		helloWorldCID = "bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea"
		dbdir         = "db"
		workdir       = "repo"
		erase         = true
	)

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	db, err := pebble.Open(dbdir, &pebble.Options{})
	if err != nil {
		return fmt.Errorf("could not open pebble db: %w", err)
	}

	store := store.New(db, codec.NewJSONCodec())

	fstore := fstore.New(log, store, workdir)

	log.Info().Msg("created fstore")

	defer func(erase bool) {
		if !erase {
			return
		}

		os.RemoveAll(dbdir)
		os.RemoveAll(workdir)
	}(erase)

	cids := []string{
		helloWorldCID,
		counterCID,
	}

	for _, cid := range cids {

		err = fstore.Install(manifestURLFromCID(cid), cid)
		if err != nil {
			return fmt.Errorf("could not install function: (%s): %w", cid, err)
		}

	}

	return nil
}

func manifestURLFromCID(cid string) string {
	return fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid)
}
