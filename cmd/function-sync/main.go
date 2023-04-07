package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/fstore"
	"github.com/blocklessnetworking/b7s/store"
)

const (
	success = 0
	failure = 1
)

var (
	log zerolog.Logger
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagDownload  bool
		flagDB        string
		flagWorkspace string
	)

	pflag.StringVarP(&flagDB, "database", "d", "function-db", "function database to use")
	pflag.StringVarP(&flagWorkspace, "workspace", "w", "workspace", "workspace to use")
	pflag.BoolVar(&flagDownload, "download", false, "download files only, don't sync")

	pflag.CommandLine.SortFlags = false

	pflag.Parse()

	log = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Open the pebble function database.
	fdb, err := pebble.Open(flagDB, &pebble.Options{})
	if err != nil {
		log.Error().Err(err).Str("db", flagDB).Msg("could not open pebble function database")
		return failure
	}
	defer fdb.Close()

	store := store.New(fdb)

	// Create function store.
	fstore := fstore.New(log, store, flagWorkspace)

	if flagDownload {
		return download(fstore)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	runSyncLoop(ctx)

	return success
}

func download(fstore *fstore.FStore) int {

	functions := []struct {
		Address string
		CID     string
	}{
		{
			Address: "https://bafybeie67hbrd5bbl6blkhu5xs7mn2aq53ngh2m6vcmyjwrjbvsmh2scqu.ipfs.w3s.link/manifest.json",
			CID:     "bafybeie67hbrd5bbl6blkhu5xs7mn2aq53ngh2m6vcmyjwrjbvsmh2scqu",
		},
		{
			Address: "https://bafybeig32qk2p4mrphle4mjix4lctijnfujemrc5l4lmusbmyi6qpaw3va.ipfs.w3s.link/manifest.json",
			CID:     "bafybeig32qk2p4mrphle4mjix4lctijnfujemrc5l4lmusbmyi6qpaw3va",
		},
		{
			Address: "https://bafybeiebrkukqliurwdwlinnpocpkk62qz462enez5adfcyp2tgy3aa2mq.ipfs.w3s.link/manifest.json",
			CID:     "bafybeiebrkukqliurwdwlinnpocpkk62qz462enez5adfcyp2tgy3aa2mq",
		},
		{
			Address: "https://bafybeiguusihksjk4uln2eeaom7vcgmrwcngfyfeli63uiem4m4zzdz4pu.ipfs.w3s.link/manifest.json",
			CID:     "bafybeiguusihksjk4uln2eeaom7vcgmrwcngfyfeli63uiem4m4zzdz4pu",
		},
		{
			Address: "https://bafybeih7sjfdm2rrcpgs3o44kvx3l2ec3svvnsmy35adyxr4dvvqptlik4.ipfs.w3s.link/manifest.json",
			CID:     "bafybeih7sjfdm2rrcpgs3o44kvx3l2ec3svvnsmy35adyxr4dvvqptlik4",
		},
		{
			Address: "https://bafybeifzg5ge3dt3kitzgaedk2ji4h3mhsmzd3kp3mrstkwzkruo2wp7n4.ipfs.w3s.link/manifest.json",
			CID:     "bafybeifzg5ge3dt3kitzgaedk2ji4h3mhsmzd3kp3mrstkwzkruo2wp7n4",
		},
		{
			Address: "https://bafybeifxsjz7a7zjqeeqmvwpg3nmxvkqueudgqjknlwlsn3o5rw5dv636q.ipfs.w3s.link/manifest.json",
			CID:     "bafybeifxsjz7a7zjqeeqmvwpg3nmxvkqueudgqjknlwlsn3o5rw5dv636q",
		},
		{
			Address: "https://bafybeieduydqrzxc7ils6cldwcefk75ycfgoqqmsyt42l7iejkn4rptd4m.ipfs.w3s.link/manifest.json",
			CID:     "bafybeieduydqrzxc7ils6cldwcefk75ycfgoqqmsyt42l7iejkn4rptd4m",
		},
		{
			Address: "https://bafybeibb3jytl56fagu6yczas4tz4hm7bckovwmi4xfsg5q7yrgvvthsw4.ipfs.w3s.link/manifest.json",
			CID:     "bafybeibb3jytl56fagu6yczas4tz4hm7bckovwmi4xfsg5q7yrgvvthsw4",
		},
	}

	for i, fn := range functions {

		_, err := fstore.Get(fn.Address, fn.CID, false)
		if err != nil {
			log.Error().Err(err).Str("cid", fn.CID).Msg("could not retrieve function")
			return failure
		}

		log.Info().Int("i", i).Str("cid", fn.CID).Msg("function downloaded okay")
	}

	return success
}

const (
	syncInterval = 5 * time.Second
)

func runSyncLoop(ctx context.Context) {

	ticker := time.NewTicker(syncInterval)

	for {
		select {
		case <-ticker.C:
			fmt.Printf("running sync\n")

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
