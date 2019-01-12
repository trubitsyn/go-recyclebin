// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package recyclebin

import (
	"github.com/spf13/afero"
	"net/url"
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
	if err := f.Close(); err != nil {
		t.Error("unable to close test file.")
	}
	trashedFilename := getTrashedFilename(trashPath, filename)
	if err = bin.Recycle(filename); err != nil {
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
	trashInfo, err := readTrashInfo(buildTrashInfoPath(trashPath, trashedFilename))
	if err != nil {
		t.Error("trash info ''" + trashedFilename + "'.trashinfo cannot be read")
	}
	escapedRealPath := url.PathEscape(filename)
	if escapedRealPath != trashInfo.Path {
		t.Error("trash info '" + trashedFilename + "'.trashinfo has invalid Path value")
	}
}

func TestDeleteFromTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	filename := "file"
	if err := createTrashFiles(trashPath, filename); err != nil {
		t.Error("unable to create test trash files")
	}
	if err := bin.Remove(filename); err != nil {
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
	if err := createTrashFiles(trashPath, trashFilename); err != nil {
		t.Error("unable to create test trash files")
	}
	if err := bin.Restore(trashFilename); err != nil {
		t.Error("unable to restore file '" + trashFilename + "'")
	}
	if exists, _ := afero.Exists(fs, buildTrashFilePath(trashPath, trashFilename)); exists {
		t.Error("trash still contains the file")
	}
	if exists, _ := afero.Exists(fs, buildTrashInfoPath(trashPath, trashFilename)); exists {
		t.Error("trash still contains trashinfo file")
	}
	if exists, _ := afero.Exists(fs, trashFilename); !exists {
		t.Error("file has not been restored")
	}
}

func TestEmptyTrash(t *testing.T) {
	trashPath := ".local/share/Trash"
	bin := NewRecycleBin(trashPath)
	if err := createTrashFiles(trashPath, "script.sh"); err != nil {
		t.Error("unable to create test trash files")
	}
	if err := createTrashFiles(trashPath, "lib.so"); err != nil {
		t.Error("unable to create test trash files")
	}
	if err := bin.Empty(); err != nil {
		t.Error("unable to empty the trash")
	}
	if existsTrashFile(trashPath, "script.sh") || existsTrashFile(trashPath, "lib.so") {
		t.Error("trash files were not deleted")
	}
}

func createTrashFiles(trashPath string, filename string) error {
	if err := fs.MkdirAll(trashPath+"/files", os.ModeDir); err != nil {
		return err
	}
	f, err := fs.Create(buildTrashFilePath(trashPath, filename))
	if err != nil {
		return err
	}
	_ = f.Close()
	if err := fs.MkdirAll(trashPath+"/info", os.ModeDir); err != nil {
		return err
	}
	f, err = fs.Create(buildTrashInfoPath(trashPath, filename))
	if err != nil {
		return err
	}
	defer f.Close()
	content := "[Trash Info]\n" + "Path=" + filename + "\n" + "DeletionDate=2018-10-11\n"
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func existsTrashFile(trashPath string, filename string) bool {
	hasFile, _ := afero.Exists(fs, buildTrashFilePath(trashPath, filename))
	return hasFile
}

func existsTrashInfo(trashPath string, filename string) bool {
	hasTrashInfo, _ := afero.Exists(fs, buildTrashInfoPath(trashPath, filename))
	return hasTrashInfo
}

func initHomeEnvironment() {
	os.Setenv("XDG_DATA_HOME", "/home/user")
	fs.MkdirAll("/home/user/.local/share/Trash", os.ModeDir)
}

func deinitHomeEnvironment() {
	os.Unsetenv("XDG_DATA_HOME")
}
