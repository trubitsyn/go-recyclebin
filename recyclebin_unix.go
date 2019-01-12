// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"os"
	"path"
	"time"
)

type unixRecycleBin struct {
	Path string
}

func NewRecycleBin(location string) RecycleBin {
	bin := new(unixRecycleBin)
	bin.Path = location
	return bin
}

func ForLocation(location string) (RecycleBin, error) {
	var envStorage osEnvStorage
	uid := os.Getuid()
	dir, err := getTrashDirectory(location, envStorage, uid)
	if err != nil {
		return nil, err
	}
	return NewRecycleBin(dir), nil
}

// Recycle moves file to trash.
func (bin unixRecycleBin) Recycle(filepath string) error {
	_, filename := path.Split(filepath)
	if err := fs.MkdirAll(bin.Path+"/files", os.ModeDir); err != nil {
		return err
	}
	trashedFilename := getTrashedFilename(bin.Path, filename)
	if err := fs.Rename(filepath, buildTrashFilePath(bin.Path, trashedFilename)); err != nil {
		return err
	}
	deletionDate := time.Now().Format("2006-01-02T15:04:05")
	err := writeTrashInfo(bin.Path, filepath, deletionDate, trashedFilename)
	return err
}

// Restore restores file from trash.
func (bin unixRecycleBin) Restore(trashFilename string) error {
	trashInfoPath := buildTrashInfoPath(bin.Path, trashFilename)
	trashInfo, err := readTrashInfo(trashInfoPath)
	if err != nil {
		return err
	}
	deletedFilePath := buildTrashFilePath(bin.Path, trashFilename)
	if err := fs.Rename(deletedFilePath, trashInfo.Path); err != nil {
		return err
	}
	err = fs.Remove(buildTrashInfoPath(bin.Path, trashFilename))
	return err
}

// Remove permanently deletes file from trash.
func (bin unixRecycleBin) Remove(trashFilename string) error {
	if err := fs.Remove(buildTrashFilePath(bin.Path, trashFilename)); err != nil {
		return err
	}
	err := fs.Remove(buildTrashInfoPath(bin.Path, trashFilename))
	return err
}

// Empty empties the trash.
func (bin unixRecycleBin) Empty() error {
	if err := fs.RemoveAll(bin.Path + "/files"); err != nil {
		return err
	}
	err := fs.RemoveAll(bin.Path + "/info")
	return err
}
