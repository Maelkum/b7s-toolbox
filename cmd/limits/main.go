package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sys/windows"

	"github.com/blocklessnetworking/b7s/executor/limits"
	"github.com/blocklessnetworking/b7s/models/execute"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	cpuRate := 0.5
	memLimit := int64(256 * 1000)

	log := zerolog.New(os.Stderr)

	limiter, err := limits.New(limits.WithCPUPercentage(cpuRate), limits.WithMemoryKB(memLimit))
	if err != nil {
		log.Fatal().Err(err).Msg("could not create limiter")
		return failure
	}
	defer limiter.Shutdown()

	log.Info().Msg("limiter created")

	count := 10

	var wg sync.WaitGroup

	wg.Add(count)

	for i := 0; i < count; i++ {

		go func() {

			defer wg.Done()

			// Start a process.
			sleepCmd := `C:\Program Files\Git\usr\bin\sleep.exe`
			cmd := exec.Command(sleepCmd, "30")
			err = cmd.Start()
			if err != nil {
				log.Fatal().Err(err).Msg("could not start command")
				return
			}

			log.Info().Int("pid", cmd.Process.Pid).Msg("process started")

			childHandle, err := readHandle(cmd)
			if err != nil {
				log.Fatal().Err(err).Msg("could not read process handle")
				return
			}

			// Create a duplicate handle - only for me (current process), not inheritable.
			var handle windows.Handle
			me := windows.CurrentProcess()
			err = windows.DuplicateHandle(me, childHandle, me, &handle, windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_TERMINATE|windows.PROCESS_SET_QUOTA, false, 0)
			if err != nil {
				log.Fatal().Err(err).Msg("could not duplicate process handle")
				return
			}
			defer windows.CloseHandle(handle)

			proc := execute.ProcessID{
				PID:    cmd.Process.Pid,
				Handle: uintptr(childHandle),
			}

			err = limiter.LimitProcess(proc)
			if err != nil {
				log.Fatal().Err(err).Msg("could not set limits for process")
			}

			err = cmd.Wait()
			if err != nil {
				log.Fatal().Err(err).Msg("could not wait command")
				return
			}

			log.Info().Msg("external command completed")
		}()
	}

	time.Sleep(1 * time.Second)

	log.Info().Msg("process limits set")

	pids, err := limiter.ListProcesses()
	if err != nil {
		log.Error().Err(err).Msg("could not get list of limited processes")
	} else {
		log.Info().
			Int("count", len(pids)).
			Msg("received list of limited processes")

		for i, pid := range pids {
			fmt.Printf("pid %v: %v\n", i, pid)
		}
	}

	wg.Wait()

	log.Info().Msg("all done")

	return success
}

func readHandle(cmd *exec.Cmd) (windows.Handle, error) {

	proc := cmd.Process
	if proc == nil {
		return 0, fmt.Errorf("command not started")
	}

	v := reflect.ValueOf(proc).Elem()
	field := v.FieldByName("handle")

	if field.IsZero() {
		return 0, fmt.Errorf("field not found")
	}

	// NOTE: Returning uintptr as uint64.
	handle := windows.Handle(field.Uint())

	return handle, nil
}
