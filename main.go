package main

import (
	"bomoko/lagoon-init/cmd"
	"os"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
