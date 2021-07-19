// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSizes(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}
	t.Log("test how different arguments affect binary and archive sizes")
	cases := []struct {
		file string
		args []string
	}{
		{"testdata/apigw.go", nil},
		{"testdata/noop.go", nil},
		{"testdata/noop.go", []string{"-tags", "lambda.norpc"}},
		{"testdata/noop.go", []string{"-ldflags=-s -w"}},
		{"testdata/noop.go", []string{"-tags", "lambda.norpc", "-ldflags=-s -w"}},
	}
	testDir, err := os.Getwd()
	require.NoError(t, err)
	tempDir, err := ioutil.TempDir("/tmp", "build-lambda-zip")
	require.NoError(t, err)
	for _, test := range cases {
		require.NoError(t, os.Chdir(testDir))
		testName := fmt.Sprintf("%s, %v", test.file, test.args)
		t.Run(testName, func(t *testing.T) {
			binPath := path.Join(tempDir, test.file+".bin")
			zipPath := path.Join(tempDir, test.file+".zip")

			buildArgs := []string{"build", "-o", binPath}
			buildArgs = append(buildArgs, test.args...)
			buildArgs = append(buildArgs, test.file)

			gocmd := exec.Command("go", buildArgs...)
			gocmd.Env = append(os.Environ(), "GOOS=linux")
			gocmd.Stderr = os.Stderr
			require.NoError(t, gocmd.Run())
			require.NoError(t, os.Chdir(filepath.Dir(binPath)))
			require.NoError(t, compressExeAndArgs(zipPath, binPath, []string{}))

			binInfo, err := os.Stat(binPath)
			require.NoError(t, err)
			zipInfo, err := os.Stat(zipPath)
			require.NoError(t, err)

			t.Logf("zip size = %d Kb, bin size = %d Kb", zipInfo.Size()/1024, binInfo.Size()/1024)
		})
	}

}

func TestCompressExeAndArgs(t *testing.T) {
	tempDir, err := ioutil.TempDir("/tmp", "build-lambda-zip")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fileNames := []string{"the-exe", "the-config", "other-config"}
	filePaths := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		filePaths[i] = filepath.Join(tempDir, fileName)
		f, err := os.Create(filePaths[i])
		require.NoError(t, err)
		_, _ = fmt.Fprintf(f, "Hello file %d!", i)
		err = f.Close()
		require.NoError(t, err)
	}
	outZipPath := filepath.Join(tempDir, "lambda.zip")

	err = compressExeAndArgs(outZipPath, filePaths[0], filePaths[1:])
	require.NoError(t, err)

	t.Run("handler exe configured in zip root", func(t *testing.T) {
		require.NotEqual(t, filePaths[0], filepath.Base(filePaths[0]), "test precondition")
		zipReader, err := zip.OpenReader(outZipPath)
		require.NoError(t, err)
		defer zipReader.Close()
		for _, zf := range zipReader.File {
			if zf.Name == filepath.Base(filePaths[0]) {
				assert.True(t, zf.FileInfo().Mode().IsRegular())
				permissions := int(zf.FileInfo().Mode() & 0777)
				assert.Equal(t, 0111, permissions&0111, "file permissions: %#o, aren't executable", permissions)
				assert.Equal(t, 0444, permissions&0444, "file permissions: %#o, aren't readable", permissions)
				return
			}
		}
		t.Fatalf("failed to find handler exe in zip")
	})

	t.Run("boostrap is a symlink to handler exe", func(t *testing.T) {
		zipReader, err := zip.OpenReader(outZipPath)
		require.NoError(t, err)
		defer zipReader.Close()
		var bootstrap *zip.File
		for _, f := range zipReader.File {
			if f.Name == "bootstrap" {
				bootstrap = f
			}
		}
		require.NotNil(t, bootstrap)
		assert.Equal(t, 0755|os.ModeSymlink, bootstrap.FileInfo().Mode())
		link, err := bootstrap.Open()
		require.NoError(t, err)
		defer link.Close()
		linkTarget, err := ioutil.ReadAll(link)
		require.NoError(t, err)
		assert.Equal(t, filepath.Base(filePaths[0]), string(linkTarget))
	})

	t.Run("resource file paths", func(t *testing.T) {
		zipReader, err := zip.OpenReader(outZipPath)
		require.NoError(t, err)
		defer zipReader.Close()
	eachFile:
		for _, path := range filePaths[1:] {
			for _, zipFileEntry := range zipReader.File {
				if zipFileEntry.Name == path {
					continue eachFile
				}
			}
			t.Logf("failed to find resource file %s in zip", path)
			t.Fail()
		}
	})

	t.Run("file contents match", func(t *testing.T) {
		zipReader, err := zip.OpenReader(outZipPath)
		require.NoError(t, err)
		defer zipReader.Close()
		expectedIndex := 0
		for _, zf := range zipReader.File {
			if zf.FileInfo().Mode().IsRegular() {
				f, err := zf.Open()
				require.NoError(t, err)
				defer f.Close()
				content, err := ioutil.ReadAll(f)
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Hello file %d!", expectedIndex), string(content), "in file: %s", zf.Name)
				expectedIndex++
			}
		}
	})

}
