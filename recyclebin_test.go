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
	bin, err := ForLocation(trashPath)
	if err != nil {
		t.Error("unable to create recycle bin.")
	}
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
	bin, err := ForLocation(trashPath)
	if err != nil {
		t.Error("unable to create recycle bin.")
	}
	filename := "file"
	createTrashFile(trashPath, filename)
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

func TestRestoreFromTrash(t *testing.T) {}

func TestEmptyTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin, err := ForLocation(trashPath)
	if err != nil {
		t.Error("unable to create recycle bin.")
	}
	createTrashFile(trashPath, "script.sh")
	createTrashFile(trashPath, "lib.so")
	err = bin.Empty()
	success := err == nil && !existsTrashFile(trashPath, "script.sh") && !existsTrashFile(trashPath, "lib.so")
	if !success {
		t.Error("trash has not been emptied")
	}
}

func TestEmptyHomeTrash(t *testing.T) {
	// TODO
}

func TestEmptyDeviceTrash(t *testing.T) {
	// TODO
}

func TestMoveToHomeTrash(t *testing.T) {
	// TODO
}

func TestMoveToDeviceTrash(t *testing.T) {
	// TODO
}

func TestRemoveFromHomeTrash(t *testing.T) {
	// TODO
}

func TestRemoveFromDeviceTrash(t *testing.T) {
	// TODO
}

func TestRestoreFromHomeTrash(t *testing.T) {
	// TODO
}

func TestRestoreFromDeviceTrash(t *testing.T) {
	// TODO
}

func createTrashFile(trashPath string, filename string) {
	fs.MkdirAll(trashPath+"/files", os.ModeDir)
	fs.Create(trashPath + "/files/" + filename)
	fs.MkdirAll(trashPath+"/info", os.ModeDir)
	fs.Create(trashPath + "/info/" + filename + ".trashinfo")
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
