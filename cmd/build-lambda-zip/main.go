// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const usage = `build-lambda-zip - Puts an executable and supplemental files into a zip file that works with AWS Lambda.
usage:
  build-lambda-zip [options] handler-exe [paths...]
options:
  -o, --output  output file path for the zip. (default: ${handler-exe}.zip)
  -h, --help    prints usage
`

func main() {
	var outputZip string
	flag.StringVar(&outputZip, "o", "", "")
	flag.StringVar(&outputZip, "output", "", "")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no input provided")
	}
	inputExe := flag.Arg(0)
	if outputZip == "" {
		outputZip = fmt.Sprintf("%s.zip", filepath.Base(inputExe))
	}
	if err := compressExeAndArgs(outputZip, inputExe, flag.Args()[1:]); err != nil {
		log.Fatalf("failed to compress file: %v", err)
	}
	log.Printf("wrote %s", outputZip)
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
	if pathInZip != "bootstrap" {
		header := &zip.FileHeader{Name: "bootstrap", Method: zip.Deflate}
		header.SetMode(0755 | os.ModeSymlink)
		link, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err := link.Write([]byte(pathInZip)); err != nil {
			return err
		}
	}

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
	return err
}

func compressExeAndArgs(outZipPath string, exePath string, args []string) error {
	zipFile, err := os.Create(outZipPath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := zipFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close zip file: %v\n", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		return err
	}

	err = writeExe(zipWriter, filepath.Base(exePath), data)
	if err != nil {
		return err
	}

	for _, arg := range args {
		writer, err := zipWriter.Create(arg)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	return err
}
