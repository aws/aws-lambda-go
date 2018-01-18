package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	err := compressExe(fmt.Sprintf("%s.zip", path.Base(os.Args[1])), os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
	log.Print(pathInZip)
	exe, err := writer.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           pathInZip,
		Method:         zip.Deflate,
	})
	if err != nil {
		return err
	}

	_, err = exe.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func compressExe(outZipPath string, exePath string) (err error) {
	zipFile, err := os.Create(outZipPath)
	if err != nil {
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		return
	}

	return writeExe(zipWriter, filepath.Base(exePath), data)
}
