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
	Path         string
	DeletionDate string
}

func readTrashInfo(trashInfoFile string) (TrashInfo, error) {
	file, err := fs.Open(trashInfoFile)
	if err != nil {
		return TrashInfo{}, err
	}

	reader := bufio.NewReader(file)
	header, _, _ := reader.ReadLine()
	pathPair, _, _ := reader.ReadLine()
	deletionDatePair, _, _ := reader.ReadLine()

	if string(header) != "[Trash Info]" {
		return TrashInfo{}, errors.New(".trashinfo file is not valid")
	}

	path := strings.Split(string(pathPair), "=")[1]
	deletionDate := strings.Split(string(deletionDatePair), "=")[1]

	return TrashInfo{path, deletionDate}, nil
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
