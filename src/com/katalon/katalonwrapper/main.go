package main

import (
	"com/katalon/katalonwrapper/download"
	"fmt"
)

func main() {
	//ksVersion := os.Args[1]
	//download.GetVersion(ksVersion)
	var version = "6.3.3"
	katalonDir := download.GetKatalonPackage(version)
	fmt.Println("Katalon Directory:", katalonDir)
}
