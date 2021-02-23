package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sharpevo/seqbot/cmd/send/app"
)

func main() {
	sendFlagSet := flag.NewFlagSet("send", flag.ExitOnError)
	cmd := app.NewSendCommand(sendFlagSet)
	sendFlagSet.Parse(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
