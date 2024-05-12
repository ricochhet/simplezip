package simplezip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ZipMessenger struct {
	AddedFile func(string)
}

func DefaultZipMessenger() ZipMessenger {
	return ZipMessenger{
		AddedFile: func(path string) {
			fmt.Println("Adding file to zip: " + path)
		},
	}
}

func Zip(dirPath string, outPath string) error {
	return ZipWithMessenger(dirPath, outPath, DefaultZipMessenger())
}

func ZipWithMessenger(dirPath string, outPath string, messenger ZipMessenger) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o700); err != nil {
		return err
	}

	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWrite := zip.NewWriter(file)
	defer zipWrite.Close()

	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		zipCreate, err := zipWrite.Create(convertPath(path, dirPath))
		if err != nil {
			return err
		}

		messenger.AddedFile(path)

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		if _, err = io.Copy(zipCreate, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func convertPath(path, src string) string {
	path = trimSrcPrefix(path, src)
	path = replaceBackslashes(path)

	return path
}

func trimSrcPrefix(path, src string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, src), string(filepath.Separator))
}

func replaceBackslashes(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
