package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const dftNTPServer = "pool.ntp.org"

func main() {
	t, err := ntp.Time(dftNTPServer)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error getting time from %s: %v\n", dftNTPServer, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Time: %v\n", t)
}
