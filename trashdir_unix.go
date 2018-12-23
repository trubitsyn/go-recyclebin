// +build !windows

package recyclebin

import (
	"errors"
	"github.com/spf13/afero"
	"os"
	"strconv"
)

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
