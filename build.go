package main

import (
	"strings"
	"os/exec"
	"os"
	"archive/zip"
	"path/filepath"
	"io"
)

func main() {
	buildGo("measurePerformance/main.go")
}

func buildGo(goHandler string) {

	executable := strings.TrimSuffix(goHandler, ".go")
	cmd := exec.Command("go", "build", "-i", "-o", executable, goHandler)
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")

	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	err = zipFile(executable + ".zip", executable)
	if err != nil {
		panic(err)
	}
}

func zipFile(zipFile string, file string) error {

	f, err := os.Create(zipFile)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(f)

	fileForZip, _ := os.Open(file)

	writer, _ := zipWriter.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           filepath.Base(file),
		Method:         zip.Deflate,
	})

	io.Copy(writer, fileForZip)
	fileForZip.Close()
	zipWriter.Close()

	return nil
}
