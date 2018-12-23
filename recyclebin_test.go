// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package recyclebin

import (
	"github.com/spf13/afero"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
	fs = afero.NewMemMapFs()
}

func teardown() {
	os.Unsetenv("XDG_DATA_HOME")
}

func TestMoveToTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	filename := "file"
	f, err := fs.Create(filename)
	if err != nil {
		t.Error("unable to create test file for removal.")
	}
	err = f.Close()
	if err != nil {
		t.Error("unable to close test file.")
	}
	trashedFilename := getTrashedFilename(trashPath, filename)
	err = bin.Recycle(filename)
	if err != nil {
		t.Error("unable to recycle test file.")
	}
	fileExists, err := afero.Exists(fs, filename)
	if err != nil {
		t.Error("unable to check if file is still not deleted")
	}
	if fileExists {
		t.Error("file has not been moved to trash")
	}
	if !existsTrashFile(trashPath, trashedFilename) {
		t.Error("trash file '" + trashedFilename + "' has not been created")
	}
	if !existsTrashInfo(trashPath, trashedFilename) {
		t.Error("trash info '" + trashedFilename + ".trashinfo' has not been created")
	}
}

func TestDeleteFromTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	filename := "file"
	err := createTrashFiles(trashPath, filename)
	if err != nil {
		t.Error("unable to create test trash files")
	}
	err = bin.Remove(filename)
	if err != nil {
		t.Error("unable to remove file")
	}
	if existsTrashFile(trashPath, filename) {
		t.Error("trash file '" + filename + "' has not been removed")
	}
	if existsTrashInfo(trashPath, filename) {
		t.Error("trash info '" + filename + ".trashinfo' has not been removed")
	}
}

func TestRestoreFromTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	trashFilename := "file"
	err := createTrashFiles(trashPath, trashFilename)
	if err != nil {
		t.Error("unable to create test trash files")
	}
	err = bin.Restore(trashFilename)
	if err != nil {
		t.Error("unable to restore file '" + trashFilename + "'")
	}
	if exists, _ := afero.Exists(fs, trashPath+"/files/"+trashFilename); exists {
		t.Error("trash still contains the file")
	}
	if exists, _ := afero.Exists(fs, trashPath+"/info/"+trashFilename+".trashinfo"); exists {
		t.Error("trash still contains trashinfo file")
	}
	if exists, _ := afero.Exists(fs, trashFilename); !exists {
		t.Error("file has not been restored")
	}
}

func TestEmptyTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	err := createTrashFiles(trashPath, "script.sh")
	if err != nil {
		t.Error("unable to create test trash files")
	}
	err = createTrashFiles(trashPath, "lib.so")
	if err != nil {
		t.Error("unable to create test trash files")
	}
	err = bin.Empty()
	if err != nil {
		t.Error("unable to empty the trash")
	}
	if existsTrashFile(trashPath, "script.sh") || existsTrashFile(trashPath, "lib.so") {
		t.Error("trash files were not deleted")
	}
}

func createTrashFiles(trashPath string, filename string) error {
	err := fs.MkdirAll(trashPath+"/files", os.ModeDir)
	if err != nil {
		return err
	}
	f, err := fs.Create(trashPath + "/files/" + filename)
	if err != nil {
		return err
	}
	_ = f.Close()
	err = fs.MkdirAll(trashPath+"/info", os.ModeDir)
	if err != nil {
		return err
	}
	f, err = fs.Create(trashPath + "/info/" + filename + ".trashinfo")
	if err != nil {
		return err
	}
	f.WriteString("[Trash Info]\n")
	f.WriteString("Path=" + filename + "\n")
	f.WriteString("DeletionDate=2018-10-11\n")
	_ = f.Close()
	return nil
}

func existsTrashFile(trashPath string, filename string) bool {
	hasFile, _ := afero.Exists(fs, trashPath+"/files/"+filename)
	return hasFile
}

func existsTrashInfo(trashPath string, filename string) bool {
	hasTrashInfo, _ := afero.Exists(fs, trashPath+"/info/"+filename+".trashinfo")
	return hasTrashInfo
}

func initHomeEnvironment() {
	os.Setenv("XDG_DATA_HOME", "/home/user")
	fs.MkdirAll("/home/user/.local/share/Trash", os.ModeDir)
}

func deinitHomeEnvironment() {
	os.Unsetenv("XDG_DATA_HOME")
}
