package main

import (
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/cmd/NEXTSEQ-500/app"
)

func main() {
	cmd := app.NewNextseq500Command()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
