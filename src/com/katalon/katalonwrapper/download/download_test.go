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
