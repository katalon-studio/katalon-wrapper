package main

import (
	"com/katalon/katalonwrapper/download"
	"com/katalon/katalonwrapper/executor"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	BuildVersion = ""
)

func commandLineUsage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] katalon-version [argument...]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var ksVersion, proxyURL string
	var printVersion bool

	flag.Usage = commandLineUsage

	flag.StringVar(&proxyURL, "proxy", "", "Proxy server address (i.e. http://[host]:[port])")
	flag.BoolVar(&printVersion, "version", false, "Current version")

	flag.Parse()

	if printVersion {
		fmt.Println("Katalon Wrapper version:", BuildVersion)
		return
	}

	remainingArgs := flag.Args()
	nArgs := flag.NArg()
	if nArgs < 1 {
		log.Fatal("Katalon version must be provided.")
	}

	log.Println("Katalon Wrapper version:", BuildVersion)

	ksVersion = remainingArgs[0]

	katalonDir := download.GetKatalonPackage(ksVersion, proxyURL)

	if nArgs > 1 {
		// Execute Katalon command
		commandArgs := remainingArgs[1:]
		executor.Execute(katalonDir, commandArgs)
	} else {
		log.Println("Arguments are empty, no test will be executed.")
	}
}
