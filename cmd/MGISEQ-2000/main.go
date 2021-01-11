package main

import (
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/cmd/MGISEQ-2000/app"
)

func main() {
	cmd := app.NewMgi2000Command()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
