package download

import (
	"net/http"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

const Releases = "https://raw.githubusercontent.com/katalon-studio/katalon-studio/master/releases.json"

type KatalonVersion struct{
	Version, Filename, Os, Url, ContainingFolder string
}

func GetOS() string{
	os := runtime.GOOS
	arch := runtime.GOARCH
	if os == "windows" {
		if arch == "amd64" {
			return "windows 64"
		} else {
			return "windows 32"
		}
	} else if os == "darwin" {
		return "macos (app)"
	} else if os == "linux" {
		return "linux"
	}
	return ""
}

func GetVersion(ksVersion string) (version KatalonVersion) {
	resp, err := http.Get(Releases)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	katalonVersions := make([]KatalonVersion,0)
	err = json.NewDecoder(resp.Body).Decode(&katalonVersions)
	if err != nil {
		panic(err)
	}
	os := GetOS()
	for _, v := range katalonVersions {
		if v.Version == ksVersion && strings.EqualFold(v.Os, os) {
			version = v
			fmt.Println(version)
			return
		}
	}
	return
}