package download

import (
	"strings"
	"testing"
)

const version = "6.3.3"

func TestGetKatalonPackage(t *testing.T) {
	katalonDir := GetKatalonPackage(version, "")
	t.Log("Katalon Directory", katalonDir)
}

func TestGetVersion(t *testing.T) {
	katalonVersion := GetVersion(version)
	if katalonVersion.Version == version && strings.Contains(katalonVersion.Filename, "6.3.3") {
		t.Logf("GetVersion(%s)=%v", version, katalonVersion)
	} else {
		t.Errorf("Invalid katalon version %v", katalonVersion)
	}
}

func TestExtractFileZip(t *testing.T) {
	macZip := "Katalon.Studio.app.zip"
	extractedFiles, err := ExtractFile(macZip, ".")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(extractedFiles)
}

func TestExtractFileTar(t *testing.T) {
	linuxTar := "Katalon_Studio_Linux_64-6.3.3.tar.gz"
	extractedFiles, err := ExtractFile(linuxTar, ".")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(extractedFiles)
}
