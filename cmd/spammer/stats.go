package main

import (
	"time"

	"github.com/blocklessnetwork/b7s/models/response"
)

var stats testStats

type testStats struct {
	requests []requestStat
}

type requestStat struct {
	id        uint64
	ts        time.Time
	respondts time.Time
}

func processPong(pong response.Pong) {
	log.Debug().Uint64("id", pong.ID).Msg("processing pong")
}

func printStats() {}
