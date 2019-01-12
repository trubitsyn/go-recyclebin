// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"bufio"
	"errors"
	"net/url"
	"strings"
)

type TrashInfo struct {
	TrashInfoMtime int64
	Path           string
	DeletionDate   string
}

func readTrashInfo(trashInfoPath string) (TrashInfo, error) {
	file, err := fs.Open(trashInfoPath)
	if err != nil {
		return TrashInfo{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	header := scanner.Text()

	scanner.Scan()
	pathPair := scanner.Text()

	scanner.Scan()
	deletionDatePair := scanner.Text()

	if string(header) != "[Trash Info]" {
		return TrashInfo{}, errors.New(".trashinfo file is not valid")
	}

	path := strings.Split(string(pathPair), "=")[1]
	pathUnescaped, err := url.PathUnescape(path)
	if err != nil {
		return TrashInfo{}, err
	}
	deletionDate := strings.Split(string(deletionDatePair), "=")[1]
	info, err := fs.Stat(trashInfoPath)
	trashInfoMtime := info.ModTime().Unix()

	return TrashInfo{trashInfoMtime, pathUnescaped, deletionDate}, nil
}

func writeTrashInfo(trashPath string, filepath string, deletionDate, trashedFilename string) error {
	f, err := fs.Create(buildTrashInfoPath(trashPath, trashedFilename))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString("[Trash Info]\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("Path=" + url.PathEscape(filepath) + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("DeletionDate=" + deletionDate + "\n"); err != nil {
		return err
	}
	return nil
}

func buildTrashInfoPath(trashPath string, filename string) string {
	return trashPath + "/info/" + filename + ".trashinfo"
}
