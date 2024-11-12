package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
)

func main() {

	var (
		flagCount         uint
		flagSleepDuration time.Duration
	)

	pflag.UintVar(&flagCount, "count", 60*10, "how many loops will the program do")
	pflag.DurationVar(&flagSleepDuration, "sleep-duration", time.Second, "pause between loops")

	pflag.Parse()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Printf("received a SIGINT!")
		os.Exit(1)
	}()

	for i := uint(0); i < flagCount-1; i++ {
		fmt.Printf("stdout #%v\n", i)
		log.Printf("stderr #%v\n", i)

		time.Sleep(flagSleepDuration)
	}
}
