package main

import (
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/models/response"
)

const (
	delimiter  = `__________________________________________________`
	timeLayout = "15:04:05.000"
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

func printStats(start, end time.Time, count uint64, rps float64) {
	fmt.Println(delimiter)

	fmt.Printf("Start time: %s\n", start.Format(timeLayout))
	fmt.Printf("End time  : %s\n", end.Format(timeLayout))

	total := end.Sub(start)
	avg := float64(total.Milliseconds()) / float64(count)

	fmt.Printf("Count: %v\n", count)
	fmt.Printf("RPS: %v\n", rps)

	fmt.Printf("Total time: %s\n", total.String())
	fmt.Printf("Time per request: %v ms\n", avg)

	fmt.Println(delimiter)
}
