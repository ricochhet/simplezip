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
