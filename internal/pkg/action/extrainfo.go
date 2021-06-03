package action

import (
	"fmt"
)

const (
	NAME_EXTRA         = "Extra"
	MSG_TPL_EXTRA_SUCC = "- Misc: %s"
	MSG_EXTRA_FAIL     = "- Misc: -"
)

type ExtraAction struct{}

func (e *ExtraAction) Run(
	eventName string,
	command CommandInterface,
) (string, error) {
	info, err := command.Sequencer().GetExtraExperimentInfo(eventName)
	if err != nil {
		return MSG_EXTRA_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_EXTRA_SUCC, info), nil
}

func (s *ExtraAction) Name() string {
	return NAME_EXTRA
}
