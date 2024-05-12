package simplezip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Zip(dirPath string, outPath string) error {
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipwrite := zip.NewWriter(file)
	defer zipwrite.Close()

	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		zipcreate, err := zipwrite.Create(convertPath(path, dirPath))
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		if _, err = io.Copy(zipcreate, file); err != nil {
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
