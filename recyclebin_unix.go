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
	trashInfo, err := readTrashInfo(bin.Path + "/info/" + trashInfoFile)
	if err != nil {
		return err
	}
	deletedFilePath := bin.Path + "/files/" + trashFilename
	err = fs.Rename(deletedFilePath, trashInfo.Path)
	if err != nil {
		return err
	}
	err = fs.Remove(bin.Path + "/info/" + trashInfoFile)
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
