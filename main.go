package main

import (
	"github.com/nveeser/vyconfigure/cmd"
	"log"
	"os"
)

func main() {
	rootCmd := cmd.NewRoot()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
