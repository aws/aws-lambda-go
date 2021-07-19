// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "build-lambda-zip",
		Usage: "Put an executable and supplemental files into a zip file that works with AWS Lambda.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "output file path for the zip. Defaults to the first input file name.",
			},
		},
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				return errors.New("no input provided")
			}

			inputExe := c.Args().First()
			outputZip := c.String("output")
			if outputZip == "" {
				outputZip = fmt.Sprintf("%s.zip", filepath.Base(inputExe))
			}

			if err := compressExeAndArgs(outputZip, inputExe, c.Args().Tail()); err != nil {
				return fmt.Errorf("failed to compress file: %v", err)
			}
			log.Print("wrote " + outputZip)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
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
