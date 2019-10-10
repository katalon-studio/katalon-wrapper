package main

import (
	"com/katalon/katalonwrapper/download"
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

	if flag.NArg() < 1 {
		log.Fatal("Katalon version must be provided.")
	}
	ksVersion = flag.Args()[0]

	download.GetKatalonPackage(ksVersion, proxyURL)
}
