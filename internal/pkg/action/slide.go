package action

import (
	"fmt"
)

const (
	NAME_SLIDE         = "Slide"
	MSG_TPL_SLIDE_SUCC = "- Slide: %s"
)

type SlideAction struct{}

func (s *SlideAction) Run(
	eventName string,
	wfqLogPath string,
	chipId string,
) (string, error) {
	return fmt.Sprintf(MSG_TPL_SLIDE_SUCC, chipId), nil
}

func (s *SlideAction) Name() string {
	return NAME_SLIDE
}
