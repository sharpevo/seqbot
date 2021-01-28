package action

import (
	"github.com/sharpevo/seqbot/internal/pkg/sequencer"
)

type CommandInterface interface {
	Execute() error
	Sequencer() sequencer.SequencerInterface
}
