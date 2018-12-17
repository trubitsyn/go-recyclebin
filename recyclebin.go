// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to use Trash (or Recycle Bin).
package recyclebin

import (
	"bufio"
	"errors"
	"github.com/spf13/afero"
	"os"
	"path"
	fpath "path/filepath"
	"strconv"
	"strings"
	"time"
)

var fs = afero.NewOsFs()

// MoveToTrash moves file to trash.
func MoveToTrash(filepath string) error {
	trashPath, err := getTrashDirectory(filepath)
	if err != nil {
		return err
	}
	_, filename := path.Split(filepath)
	fs.MkdirAll(trashPath+"/files", os.ModeDir)
	trashedFilename := getTrashedFilename(trashPath, filename)
	err = fs.Rename(filepath, trashPath+"/files/"+trashedFilename)
	if err != nil {
		return err
	}
	err = writeTrashInfo(trashPath, filepath, trashedFilename)
	return err
}

func getTrashedFilename(trashPath string, filename string) string {
	trashedFilename := filename
	isDuplicateFilename, _ := afero.Exists(fs, trashPath+"/files/"+trashedFilename)
	if isDuplicateFilename {
		trashedFilename = generateNewFilename(trashedFilename)
	}
	return trashedFilename
}

// RestoreFromTrash restores file from trash.
func RestoreFromTrash(filename string) error {
	trashInfoFile := filename + ".trashinfo"
	trashInfo, err := readTrashInfo(trashInfoFile)
	if err != nil {
		return err
	}
	deletedFilePath := "/files/" + filename
	return fs.Rename(deletedFilePath, trashInfo.Path)
}

func readTrashInfo(trashInfoFile string) (trashInfo, error) {
	file, err := fs.Open(trashInfoFile)
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

func writeTrashInfo(trashPath string, filepath string, trashedFilename string) error {
	f, err := fs.Create(trashPath + "/info/" + trashedFilename + ".trashinfo")
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

// DeleteFromTrash permanently deletes file from trash.
func DeleteFromTrash(filename string) error {
	trashPath, err := getTrashDirectory(filename)
	if err != nil {
		return err
	}
	err = fs.Remove(trashPath + "/files/" + filename)
	if err != nil {
		return err
	}
	return fs.Remove(trashPath + "/info/" + filename + ".trashinfo")
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
	hasHomeTrash, _ := afero.DirExists(fs, homeTrashPath)
	if hasHomeTrash {
		return homeTrashPath, nil
	}
	err := fs.MkdirAll(homeTrashPath, os.ModeDir)
	return homeTrashPath, err
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
	hasTrash, _ := afero.DirExists(fs, topTrashPath)
	if !hasTrash {
		topTrashUidPath := ".Trash-" + strconv.Itoa(uid)
		err := fs.Mkdir(topTrashUidPath, os.ModeDir)
		if err != nil {
			return "", err
		}
		return topTrashUidPath, nil
	}

	if isSymlink(topTrashPath) {
		return "", errors.New("device's top .Trash directory is a symbolic link")
	}

	uidTrashPath := topTrashPath + strconv.Itoa(uid)
	hasUidTrash, _ := afero.DirExists(fs, uidTrashPath)
	if !hasUidTrash {
		err := fs.Mkdir(uidTrashPath, os.ModeDir)
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

	isDuplicateFilename, _ := afero.Exists(fs, newFilename)
	for index == -1 || isDuplicateFilename {
		index += 1
		newFilename = bareName + strconv.Itoa(index) + extension
	}
	return newFilename
}

func emptyTrash(trashPath string) {
	fs.RemoveAll(trashPath + "/files")
	fs.RemoveAll(trashPath + "/info")
}

func isSymlink(path string) bool {
	file, err := fs.Stat(path)
	return err != nil || file.Mode() != os.ModeSymlink
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
