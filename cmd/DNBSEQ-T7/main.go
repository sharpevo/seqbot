package main

import (
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/cmd/DNBSEQ-T7/app"
)

func main() {
	cmd := app.NewT7Command()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
