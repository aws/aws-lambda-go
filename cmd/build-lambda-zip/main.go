package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "build-lambda-zip"
	app.Usage = "Put an executable into a zip file that works with AWS Lambda."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "output file path for the zip. Defaults to the input file name.",
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			return errors.New("No input provided")
		}

		inputExe := c.Args().First()
		outputZip := c.String("output")
		if outputZip == "" {
			outputZip = fmt.Sprintf("%s.zip", filepath.Base(inputExe))
		}

		if err := compressExe(outputZip, inputExe); err != nil {
			return fmt.Errorf("Failed to compress file: %v", err)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
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

func compressExe(outZipPath, exePath string) error {
	zipFile, err := os.Create(outZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		return err
	}

	return writeExe(zipWriter, filepath.Base(exePath), data)
}
