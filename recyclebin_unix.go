// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !windows

package recyclebin

import (
	"os"
	"path"
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
	dir, err := getTrashDirectory(location)
	if err != nil {
		return nil, err
	}
	return NewRecycleBin(dir), nil
}

// Recycle moves file to trash.
func (bin unixRecycleBin) Recycle(filepath string) error {
	_, filename := path.Split(filepath)
	err := fs.MkdirAll(bin.Path+"/files", os.ModeDir)
	if err != nil {
		return err
	}
	trashedFilename := getTrashedFilename(bin.Path, filename)
	err = fs.Rename(filepath, buildTrashFilePath(bin.Path, trashedFilename))
	if err != nil {
		return err
	}
	err = writeTrashInfo(bin.Path, filepath, trashedFilename)
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
	err = fs.Rename(deletedFilePath, trashInfo.Path)
	if err != nil {
		return err
	}
	err = fs.Remove(buildTrashInfoPath(bin.Path, trashFilename))
	return err
}

// Remove permanently deletes file from trash.
func (bin unixRecycleBin) Remove(trashFilename string) error {
	err := fs.Remove(buildTrashFilePath(bin.Path, trashFilename))
	if err != nil {
		return err
	}
	err = fs.Remove(buildTrashInfoPath(bin.Path, trashFilename))
	return err
}

// Empty empties the trash.
func (bin unixRecycleBin) Empty() error {
	err := fs.RemoveAll(bin.Path + "/files")
	if err != nil {
		return err
	}
	err = fs.RemoveAll(bin.Path + "/info")
	return err
}
