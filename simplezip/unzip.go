package simplezip

import (
	"archive/zip"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(dirPath, outPath string) error {
	return UnzipByPrefix(dirPath, outPath, "")
}

func UnzipByPrefix(dirPath, outPath, extractPrefix string) error {
	stepCopyBytes := 1024
	zipread, err := zip.OpenReader(dirPath)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, file := range zipread.File {
		readclose, err := file.Open()
		if err != nil {
			return err
		}
		defer readclose.Close()

		if extractPrefix != "" && !strings.HasPrefix(file.Name, extractPrefix) {
			continue
		}

		name, err := url.QueryUnescape(maybeTrimPrefix(file.Name, extractPrefix))
		if err != nil {
			return err
		}

		destPath := filepath.Join(outPath, name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				continue
			}
		} else if err := stepCopy(destPath, readclose, int64(stepCopyBytes)); err != nil {
			return err
		}
	}

	zipread.Close()

	return nil
}

func maybeTrimPrefix(trimmable, prefix string) string {
	if prefix != "" {
		return strings.TrimPrefix(trimmable, prefix)
	}

	return trimmable
}

func stepCopy(apath string, src io.Reader, stepCopyBytes int64) error {
	if err := os.MkdirAll(filepath.Dir(apath), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(apath)
	if err != nil {
		return err
	}

	for {
		_, err := io.CopyN(dst, src, stepCopyBytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	dst.Close()

	return nil
}
