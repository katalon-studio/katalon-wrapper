package download

import (
	"archive/tar"
	"archive/zip"
	"com/katalon/katalonwrapper/utils"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const Releases = "https://raw.githubusercontent.com/katalon-studio/katalon-studio/master/releases.json"

func getDomainName(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	hostname := u.Hostname()
	parts := strings.Split(hostname, ".")

	if len(parts) <= 2 {
		// e.g. https://github.com, http://localhost:3000
		return hostname, nil
	}

	// Note: not supporting IP address such as 127.0.0.1, so it will result in 0.0.1
	domain := strings.Join(parts[1:], ".")
	return domain, nil
}

func getHttpClient(URL, proxyURL string) (*http.Client, error) {
	var httpProxyURL func(*http.Request) (*url.URL, error) = nil
	if proxyURL != "" {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		httpProxyURL = http.ProxyURL(proxyURL)
	}

	domain, err := getDomainName(URL)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:           httpProxyURL,
		TLSClientConfig: &tls.Config{ServerName: domain},
	}

	client := &http.Client{Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		redirectDomain, err := getDomainName(req.URL.String())
		if err != nil {
			return err
		}
		// Overwrite tls config on redirect
		transport.TLSClientConfig = &tls.Config{ServerName: redirectDomain}
		return nil
	}}
	return client, nil
}

type KatalonVersion struct {
	Version, Filename, Os, Url, ContainingFolder string
}

func GetOS() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	switch osName {
	case "windows":
		if strings.Contains(arch, "64") {
			return "windows 64"
		}
		return "windows 32"
	case "darwin":
		return "macos (app)"
	case "linux":
		return "linux"
	default:
		return "linux"
	}
}

func GetVersion(ksVersion string) (version KatalonVersion) {
	resp, err := http.Get(Releases)
	utils.HandleErrorIfExists(err, "")
	defer resp.Body.Close()

	katalonVersions := make([]KatalonVersion, 0)
	err = json.NewDecoder(resp.Body).Decode(&katalonVersions)
	utils.HandleErrorIfExists(err, "")

	osInfo := GetOS()
	for _, v := range katalonVersions {
		if v.Version == ksVersion && strings.EqualFold(v.Os, osInfo) {
			version = v
			return
		}
	}
	return
}

func GetKatalonDirectory(version string) string {
	usr, _ := user.Current()
	home := usr.HomeDir

	p := filepath.Join(home, ".katalon", version)
	return p
}

func DownloadFile(fileURL string, out *os.File, proxyURL string) error {
	client, err := getHttpClient(fileURL, proxyURL)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := client.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnzipFile(src, dest string) ([]string, error) {
	var extractedPaths []string

	dest, err := filepath.Abs(dest)
	if err != nil {
		return extractedPaths, err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return extractedPaths, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		filePath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return extractedPaths, fmt.Errorf("illegal file path: %s", filePath)
		}

		extractedPaths = append(extractedPaths, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return extractedPaths, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return extractedPaths, err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return extractedPaths, err
		}

		rc, err := f.Open()
		if err != nil {
			return extractedPaths, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		_ = outFile.Close()
		_ = rc.Close()

		if err != nil {
			return extractedPaths, err
		}
	}
	return extractedPaths, nil
}

func UntarFile(src, dest string) ([]string, error) {
	var extractedPaths []string

	dest, err := filepath.Abs(dest)
	if err != nil {
		return extractedPaths, err
	}

	reader, err := os.Open(src)
	if err != nil {
		return extractedPaths, err
	}
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return extractedPaths, err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		f, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return extractedPaths, err
		}

		// Store filename/path for returning and using later on
		filePath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return extractedPaths, fmt.Errorf("illegal file path: %s", filePath)
		}

		extractedPaths = append(extractedPaths, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return extractedPaths, err
			}
			continue
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			return extractedPaths, err
		}

		_, err = io.Copy(outFile, tarReader)

		// Close the file without defer to close before next iteration of loop
		_ = outFile.Close()

		if err != nil {
			return extractedPaths, err
		}
	}
	return extractedPaths, nil
}

func ExtractFile(src, dest string) ([]string, error) {
	ext := filepath.Ext(src)
	switch ext {
	case ".zip":
		return UnzipFile(src, dest)
	case ".gz": // .tar.gz
		return UntarFile(src, dest)
	default:
		return nil, fmt.Errorf("unsupported extension %s", ext)
	}
}

func DownloadAndExtract(fileURL, targetDir, proxyURL string) {
	log.Printf("Downloading Katalon Studio from %s. It may take a few minutes.", fileURL)

	tempFile, _ := ioutil.TempFile(targetDir, "Katalon-")
	packagePath := tempFile.Name()
	defer tempFile.Close()

	err := DownloadFile(fileURL, tempFile, proxyURL)
	utils.HandleErrorIfExists(err, "Unable to download Katalon package.\n")

	log.Printf("Extract %s to %s", packagePath, targetDir)
	_, err = ExtractFile(packagePath, targetDir)

	utils.HandleErrorIfExists(err, fmt.Sprintf("Unable to extract %s to %s.\n", packagePath, targetDir))

	_ = tempFile.Close()
	_ = os.Remove(packagePath)
}

func GetKatalonPackage(version, proxyURL string) string {
	katalonDir := GetKatalonDirectory(version)
	err := utils.EnsureDir(katalonDir)
	utils.HandleErrorIfExists(err, fmt.Sprintf("Unable to create directory %s to store Katalon Studio package.\n", katalonDir))

	fileLog := filepath.Join(katalonDir, ".katalon.done")

	if utils.Exists(fileLog) {
		log.Println("Katalon Studio package has been downloaded already.")
	} else {
		err := os.RemoveAll(katalonDir)
		utils.HandleErrorIfExists(err, "")

		err = utils.EnsureDir(katalonDir)
		utils.HandleErrorIfExists(err, fmt.Sprintf("Unable to create directory %s to store Katalon Studio package.\n", katalonDir))

		katalonVersion := GetVersion(version)
		versionURL := katalonVersion.Url

		DownloadAndExtract(versionURL, katalonDir, proxyURL)
		_, err = os.Create(fileLog)
		utils.HandleErrorIfExists(err, "")
		log.Println("Katalon Studio has been installed successfully.")
	}

	katalonContainingDir := ""
	childrenNames, err := ioutil.ReadDir(katalonDir)
	utils.HandleErrorIfExists(err, "")

	for _, childrenName := range childrenNames {
		fullPath := filepath.Join(katalonDir, childrenName.Name())

		if utils.IsDirectory(fullPath) && strings.Contains(childrenName.Name(), "Katalon") {
			katalonContainingDir = fullPath
			break
		}
	}

	return katalonContainingDir
}
