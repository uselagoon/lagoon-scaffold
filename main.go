package main

import "bomoko/lagoon-init/cmd"

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}

}
