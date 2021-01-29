package action

import (
	"fmt"
)

const (
	NAME_SLIDE         = "Slide"
	MSG_TPL_SLIDE_SUCC = "- Slide: %s"
	MSG_SLIDE_FAIL     = "- Slide: -"
)

type SlideAction struct{}

func (s *SlideAction) Run(
	eventName string,
	command CommandInterface,
) (string, error) {
	slide, err := command.Sequencer().GetSlide(eventName)
	if err != nil {
		return MSG_SLIDE_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_SLIDE_SUCC, slide), nil
}

func (s *SlideAction) Name() string {
	return NAME_SLIDE
}
