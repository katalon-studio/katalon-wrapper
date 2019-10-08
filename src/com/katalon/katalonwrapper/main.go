package main

import ( 
  "com/katalon/katalonwrapper/download"
  "com/katalon/katalonwrapper/executor"
  "os"
)

func main() {
    command := os.Args[1]
    if(command == "download") {
      ksVersion := os.Args[2]
      download.GetVersion(ksVersion)
    } else if(command == "run") {
      pathToKat := os.Args[2]
      args := os.Args[3:]
      executor.Execute(pathToKat, args)
    }
}