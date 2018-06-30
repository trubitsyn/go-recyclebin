// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to use Trash (or Recycle Bin).
package recyclebin

import (
	"errors"
	"os"
	fpath "path/filepath"
	"strconv"
	"strings"
	"path"
	"bufio"
)

// MoveToTrash moves file to trash.
func MoveToTrash(filepath string) error {
	trashPath, err := getTrashDirectory(filepath)
	if err != nil {
		return err
	}
	_, filename := path.Split(filepath)
	trashedFilename := trashPath + "/files/" + filename
	if isExist(trashedFilename) {
		trashedFilename = generateNewFilename(trashedFilename)
	}
	return os.Rename(filepath, trashedFilename)
}

// RestoreFromTrash restores file from trash.
func RestoreFromTrash(filename string) error {
	trashInfoFile := filename + ".trashinfo"
	trashInfo, err := readTrashInfo(trashInfoFile)
	if err != nil {
		return err
	}
	deletedFilePath := "/files/" + filename
	return os.Rename(deletedFilePath, trashInfo.Path)
}

func readTrashInfo(trashInfoFile string) (trashInfo, error) {
	file, err := os.Open(trashInfoFile)
	if err != nil {
		return trashInfo{}, err
	}

	reader := bufio.NewReader(file)
	headerPair, _, _ := reader.ReadLine()
	pathPair, _, _ := reader.ReadLine()
	deletionDatePair, _, _ := reader.ReadLine()

	header := strings.Split(string(headerPair), "=")[1]

	if header != "[Trash Info]" {
		return trashInfo{}, errors.New(".trashinfo file is not valid")
	}

	path := strings.Split(string(pathPair), "=")[1]
	deletionDate := strings.Split(string(deletionDatePair), "=")[1]

	return trashInfo{path, deletionDate}, nil
}

// DeleteFromTrash permanently deletes file from trash.
func DeleteFromTrash(filename string) error {
	trashPath := ""
	err := os.Remove(trashPath + "/files/" + filename)
	if err != nil {
		return err
	}
	return os.Remove(trashPath + "/info/" + filename + ".trashinfo")
}

// EmptyTrash empties the trash.
func EmptyTrash() error {
	homeTrashPath, err := getHomeTrashDirectory()
	emptyTrash(homeTrashPath)
	return err
}

type trashInfo struct {
	Path         string
	DeletionDate string
}

func getTrashDirectory(filepath string) (string, error) {
	if isExternalDevice(filepath) {
		deviceTrashPath, err := getDeviceTrashDirectory(filepath)
		if err == nil {
			return deviceTrashPath, nil
		}
		return "", err
	}

	homeTrashPath, err := getHomeTrashDirectory()
	if err == nil {
		return homeTrashPath, nil
	}
	return "", errors.New("cannot find or create any trash directory")
}

func isExternalDevice(filepath string) bool {
	return false
}

func getHomeTrashDirectory() (string, error) {
	homeTrashPath := getDataHomeDirectory() + "/Trash"
	if isExist(homeTrashPath) {
		return homeTrashPath, nil
	}
	return "", errors.New("home trash directory does not exist")
}

func getDataHomeDirectory() string {
	XDG_DATA_HOME := os.Getenv("XDG_DATA_HOME")
	if XDG_DATA_HOME == "" {
		return ".local/share"
	}
	return XDG_DATA_HOME
}

func getDeviceTrashDirectory(partitionRootPath string) (string, error) {
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
		return "", errors.New("device's top .Trash directory is a symbolic link")
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

func generateNewFilename(existingFilename string) string {
	extension := fpath.Ext(existingFilename)
	bareName := strings.TrimSuffix(existingFilename, extension)
	newFilename := existingFilename
	index := -1

	for index == -1 || isExist(newFilename) {
		index += 1
		newFilename = bareName + strconv.Itoa(index) + extension
	}
	return newFilename
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

func buildTrashFilePath(trashInfoFilePath string) (string, error) {
	trashInfo, err := readTrashInfo(trashInfoFilePath)
	if err != nil {
		return "", err
	}
	return trashInfo.Path, nil
}

func buildTrashInfoPath(trashPath string, filename string) string {
	return trashPath + "/files/" + filename + ".trashinfo"
}
