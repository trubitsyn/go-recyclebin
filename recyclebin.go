// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to use Trash (or Recycle Bin).
package recyclebin

import (
	"errors"
	"os"
	"path"
	fpath "path/filepath"
	"strconv"
)

type TrashInfo struct {
	Path         string
	DeletionDate string
}

func GetTrashDirectory(filepath string) (string, error) {
	if isExternalDevice(filepath) {
		deviceTrashPath, err := GetDeviceTrashDirectory(filepath)
		if err == nil {
			return deviceTrashPath, nil
		}
		return "", err
	}

	homeTrashPath, err := GetHomeTrashDirectory()
	if err == nil {
		return homeTrashPath, nil
	}
	return "", errors.New("Cannot find or create any trash directory.")
}

func isExternalDevice(filepath string) bool {
	return false
}

func GetHomeTrashDirectory() (string, error) {
	homeTrashPath := getDataHomeDirectory() + "/Trash"
	if isExist(homeTrashPath) {
		return homeTrashPath, nil
	}
	return "", errors.New("Home trash directory does not exist.")
}

func getDataHomeDirectory() string {
	XDG_DATA_HOME := os.Getenv("XDG_DATA_HOME")
	if XDG_DATA_HOME == "" {
		return ".local/share"
	}
	return XDG_DATA_HOME
}

func GetDeviceTrashDirectory(partitionRootPath string) (string, error) {
	uid := os.Getuid()
	topTrashPath := partitionRootPath + "/.Trash"
	if !isExist(topTrashPath) {
		topTrashUidPath := ".Trash-" + strconv.Itoa(uid)
		err := os.Mkdir(topTrashUidPath, os.ModeDir)
		if err != nil {
			return "", err
		}
		return topTrashUidPath, nil
	}

	if isSymlink(topTrashPath) {
		return "", errors.New("Device top .Trash directory is a symbolic link.")
	}

	uidTrashPath := topTrashPath + strconv.Itoa(uid)
	if !isExist(uidTrashPath) {
		err := os.Mkdir(uidTrashPath, os.ModeDir)
		if err != nil {
			return "", err
		}
	}
	return uidTrashPath, nil
}

func MoveToTrash(filepath string) error {
	trashPath, err := GetTrashDirectory(filepath)
	if err != nil {
		return err
	}

	_, filename := path.Split(filepath)
	trashedFilename := trashPath + "/files/" + filename
	if isExist(trashedFilename) {
		extension := fpath.Ext(trashedFilename)
		trashedFilename = extension + "1"
	}
	return os.Rename(filepath, trashedFilename)
}

func RestoreFromTrash(filename string) {
}

func DeleteFromTrash(filename string) {
}

func EmptyTrash() {
	homeTrashPath, _ := GetHomeTrashDirectory()
	emptyTrash(homeTrashPath)
}

func emptyTrash(trashPath string) {
	os.RemoveAll(trashPath + "/files")
	os.RemoveAll(trashPath + "/info")
}

func isSymlink(path string) bool {
	file, err := os.Stat(path)
	return err != nil || file.Mode() != os.ModeSymlink
}

func isExist(path string) bool {
	dir, err := os.Stat(path)
	return err == nil && dir.Mode().IsDir()
}
