// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"bufio"
	"errors"
	"strings"
	"time"
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
	deletionDate := strings.Split(string(deletionDatePair), "=")[1]
	file.Close()
	info, err := fs.Stat(trashInfoPath)
	trashInfoMtime := info.ModTime().Unix()

	return TrashInfo{trashInfoMtime, path, deletionDate}, nil
}

func writeTrashInfo(trashPath string, filepath string, trashedFilename string) error {
	f, err := fs.Create(buildTrashInfoPath(trashPath, trashedFilename))
	if err != nil {
		return err
	}
	_, err = f.WriteString("[Trash Info]\n")
	if err != nil {
		return err
	}
	deletionDate := time.Now().Format("2006-01-02T15:04:05")
	_, err = f.WriteString("Path=" + filepath + "\n")
	if err != nil {
		return err
	}
	_, err = f.WriteString("DeletionDate=" + deletionDate + "\n")
	if err != nil {
		return err
	}
	err = f.Close()
	return err
}

func buildTrashInfoPath(trashPath string, filename string) string {
	return trashPath + "/info/" + filename + ".trashinfo"
}
