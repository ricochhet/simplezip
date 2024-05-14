/*
 * simplezip
 * Copyright (C) 2024 simplezip contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package simplezip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DefaultUnzipMessenger() ZipMessenger {
	return ZipMessenger{
		AddedFile: func(path string) {
			fmt.Println("Unzipping file: " + path)
		},
	}
}

func Unzip(dirPath, outPath string) error {
	return UnzipByPrefixWithMessenger(dirPath, outPath, "", DefaultUnzipMessenger())
}

func UnzipByPrefixWithMessenger(dirPath, outPath, extractPrefix string, messenger ZipMessenger) error {
	stepCopyBytes := 1024
	zipRead, err := zip.OpenReader(dirPath)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, file := range zipRead.File {
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
		messenger.AddedFile(destPath)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				continue
			}
		} else if err := stepCopy(destPath, readclose, int64(stepCopyBytes)); err != nil {
			return err
		}
	}

	zipRead.Close()

	return nil
}

func maybeTrimPrefix(trimmable, prefix string) string {
	if prefix != "" {
		return strings.TrimPrefix(trimmable, prefix)
	}

	return trimmable
}

func stepCopy(dirPath string, outPath io.Reader, stepCopyBytes int64) error {
	if err := os.MkdirAll(filepath.Dir(dirPath), os.ModePerm); err != nil {
		return err
	}

	destPath, err := os.Create(dirPath)
	if err != nil {
		return err
	}

	for {
		_, err := io.CopyN(destPath, outPath, stepCopyBytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	destPath.Close()

	return nil
}
