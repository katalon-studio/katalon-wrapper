package executor

import (
    "fmt"
	"os/exec"
	"strings"
	"gopkg.in/alessio/shellescape.v1"
)

func Execute(katExecutable string, args []string) {
	fmt.Println("Katalon Executable: ", katExecutable)
	
	for i, v := range args {
		splitted := strings.Split(v, "=")
		firstPart := splitted[0]
		secondPart := strings.Join(splitted[1:], "")
		args[i] = (firstPart + "=" + shellescape.Quote(secondPart))
	}
	for _, v := range args {
		fmt.Println("Katalon Argument: ", v)
	}
	RunCMD(katExecutable, args)
}

func RunCMD(path string, args []string) {
    cmd := exec.Command(path, args...)
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(string(stdout))
	cmd.Wait()
}