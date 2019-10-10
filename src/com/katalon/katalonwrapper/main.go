package main

import (
	"com/katalon/katalonwrapper/download"
	"com/katalon/katalonwrapper/executor"
	"flag"
	"fmt"
	"log"
	"os"
)

func commandLineUsage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] version\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var ksVersion, proxyURL string

	flag.Usage = commandLineUsage

	flag.StringVar(&proxyURL, "proxy", "", "Proxy server address (i.e. http://[host]:[port])")

	flag.Parse()
	flag.Usage()

	remainingArgs := flag.Args()
	nArgs := flag.NArg()
	if nArgs < 1 {
		log.Fatal("Katalon version must be provided.")
	}
	ksVersion = remainingArgs[0]

	katalonDir := download.GetKatalonPackage(ksVersion, proxyURL)

	if nArgs > 1 {
		// Execute Katalon command
		commandArgs := remainingArgs[1:]
		executor.Execute(katalonDir, commandArgs)
	}
}
