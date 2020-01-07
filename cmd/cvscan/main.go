package main

import (
	"os"
)

func main() {
	var rc int
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
		os.Exit(rc)
	}()

	cmd := newScanCmd()
	cmd.Flags().Parse(os.Args[1:])

	if err := cmd.Execute(); err != nil {
		rc = 1
	}
}
