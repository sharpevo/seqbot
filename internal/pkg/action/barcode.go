package action

import (
	"fmt"
)

const (
	NAME_BARCODE = "Barcode"

	MSG_TPL_BARCODE_SUCC = "**%s**: sequencing completed."
	MSG_BARCODE_FAIL     = "**-**: sequencing completed."
)

type BarcodeAction struct{}

func (b *BarcodeAction) Run(
	eventName string,
	command CommandInterface,
) (string, error) {
	barcode, err := command.Sequencer().GetBarcode(eventName)
	if err != nil {
		return MSG_BARCODE_FAIL, err
	}
	return fmt.Sprintf(MSG_TPL_BARCODE_SUCC, barcode), nil
}

func (b *BarcodeAction) Name() string {
	return NAME_BARCODE
}
