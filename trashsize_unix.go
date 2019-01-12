// Copyright 2019 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const directorySizesFilename = "directorysizes"

type directorySize struct {
	size          uint64
	mtime         int64
	directoryName string
}

func getDirectorySizeStruct(trashInfo TrashInfo) (directorySize, error) {
	escapedDirectoryPath := url.PathEscape(trashInfo.Path)
	trashInfoModificationTime := trashInfo.TrashInfoMtime
	dirSize, err := calculateDirectorySize(trashInfo.Path)
	if err != nil {
		return directorySize{}, err
	}

	var directorySize directorySize
	directorySize.size = dirSize
	directorySize.mtime = trashInfoModificationTime
	directorySize.directoryName = escapedDirectoryPath
	return directorySize, nil
}

func calculateDirectorySize(path string) (uint64, error) {
	var size uint64
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})
	return size, err
}

func readDirectorySizes(f *os.File) []directorySize {
	var directorySizes []directorySize
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " ")
		size, _ := strconv.ParseUint(split[0], 10, 64)
		mtime, _ := strconv.ParseInt(split[1], 10, 64)
		directoryName := split[2]

		decodedDirectoryName, err := url.PathUnescape(directoryName)
		if err != nil {
			break
		}

		directorySize := directorySize{size, mtime, decodedDirectoryName}
		directorySizes = append(directorySizes, directorySize)
	}
	return directorySizes
}

func writeDirectorySizes(f *os.File, directorySizes []directorySize) error {
	for _, directorySize := range directorySizes {
		line := string(directorySize.size) + " " + string(directorySize.mtime) + string(directorySize.directoryName) + "\n"
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func updateDirectorySizes(trashPath string, directorySizes []directorySize) error {
	tmpFilePath := trashPath + "/" + directorySizesFilename + ".tmp"
	tmpFile, err := fs.OpenFile(tmpFilePath, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	//err := writeDirectorySizes(tmpFile, directorySizes); err != nil {
	//	return err
	//}
	realFilePath := trashPath + "/" + directorySizesFilename
	err = fs.Rename(tmpFilePath, realFilePath)
	return err
}
