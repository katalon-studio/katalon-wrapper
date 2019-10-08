package main

import ( 
  "com/katalon/katalonwrapper/download"
  "os"
)

func main() {
    ksVersion := os.Args[1]
    download.GetVersion(ksVersion)
}