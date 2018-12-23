// +build !windows

package recyclebin

import (
	"bufio"
	"errors"
	"github.com/spf13/afero"
	"os"
	"path"
	"strconv"
	"strings"
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
	err = fs.Rename(filepath, bin.Path+"/files/"+trashedFilename)
	if err != nil {
		return err
	}
	err = writeTrashInfo(bin.Path, filepath, trashedFilename)
	return err
}

// Restore restores file from trash.
func (bin unixRecycleBin) Restore(trashFilename string) error {
	trashInfoFile := trashFilename + ".trashinfo"
	trashInfo, err := readTrashInfo(trashInfoFile)
	if err != nil {
		return err
	}
	deletedFilePath := "/files/" + trashFilename
	err = fs.Rename(deletedFilePath, trashInfo.Path)
	return err
}

// Remove permanently deletes file from trash.
func (bin unixRecycleBin) Remove(trashFilename string) error {
	err := fs.Remove(bin.Path + "/files/" + trashFilename)
	if err != nil {
		return err
	}
	err = fs.Remove(bin.Path + "/info/" + trashFilename + ".trashinfo")
	return err
}

// Empty empties the trash.
func (bin unixRecycleBin) Empty() error {
	err := fs.RemoveAll(bin.Path + "/files")
	if err != nil {
		return err
	}
	err = fs.RemoveAll(bin.Path + "/info")
	if err != nil {
		return err
	}
	return err
}

func getTrashDirectory(filepath string) (string, error) {
	if isExternalDevice(filepath) {
		deviceTrashPath, err := getDeviceTrashDirectory(filepath)
		if err != nil {
			return "", err
		}
		return deviceTrashPath, nil
	}

	homeTrashPath, err := getHomeTrashDirectory()
	if err != nil {
		return "", errors.New("cannot find or create any trash directory")
	}
	return homeTrashPath, nil
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
		HOME := os.Getenv("HOME")
		return HOME + "/.local/share"
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

func isSymlink(path string) bool {
	file, err := fs.Stat(path)
	return err != nil || file.Mode() != os.ModeSymlink
}

func getTrashedFilename(trashPath string, filename string) string {
	trashedFilename := filename
	isDuplicateFilename, _ := afero.Exists(fs, trashPath+"/files/"+trashedFilename)
	if isDuplicateFilename {
		trashedFilename = generateNewFilename(trashedFilename)
	}
	return trashedFilename
}

func generateNewFilename(existingFilename string) string {
	extension := path.Ext(existingFilename)
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
	headerPair, _, _ := reader.ReadLine()
	pathPair, _, _ := reader.ReadLine()
	deletionDatePair, _, _ := reader.ReadLine()

	header := strings.Split(string(headerPair), "=")[1]

	if header != "[Trash Info]" {
		return TrashInfo{}, errors.New(".trashinfo file is not valid")
	}

	path := strings.Split(string(pathPair), "=")[1]
	deletionDate := strings.Split(string(deletionDatePair), "=")[1]

	return TrashInfo{path, deletionDate}, nil
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
