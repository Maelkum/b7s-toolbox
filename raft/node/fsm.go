package node

import (
	"errors"
	"io"

	"github.com/hashicorp/raft"
	"github.com/rs/zerolog"
)

type fsm struct {
	log zerolog.Logger
}

// TODO: Implement.
func newFSM(log zerolog.Logger) *fsm {

	fsm := fsm{
		log: log,
	}

	return &fsm
}

func (f *fsm) Apply(*raft.Log) interface{} {
	f.log.Info().Msg("received Apply() call")
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.log.Info().Msg("received Snapshot() call")
	return nil, errors.New("TBD: Not implemented")
}

func (f *fsm) Restore(snapshot io.ReadCloser) error {
	f.log.Info().Msg("received Restore() call")
	return errors.New("TBD: Not implemented")
}
