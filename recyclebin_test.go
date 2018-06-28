// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package recyclebin

import (
	"testing"
	"github.com/spf13/afero"
	"os"
)

var fs = afero.NewMemMapFs()

func TestMoveToTrash(t *testing.T) {
	path := "/home/user/file"
	trashPath := "/home/user/.local/share/Trash"
	initHomeEnvironment()
	MoveToTrash("/home/user/file")

	success := !exists(path) && exists(trashPath+"/files/file") && exists(trashPath+"/info/file.trashinfo")
	if !success {
		t.Error("file has not been moved to trash")
	}
}

func TestDeleteFromTrash(t *testing.T) {
	initHomeEnvironment()
	trashPath := "/home/user/.local/share/Trash"
	fs.MkdirAll(trashPath+"/files", os.ModeDir)
	fs.Create(trashPath + "/files/file")
	fs.MkdirAll(trashPath+"/info", os.ModeDir)
	fs.Create(trashPath + "/info/file.trashinfo")
	DeleteFromTrash("file")

	success := !exists(trashPath + "/files/file") && !exists(trashPath + "/info/file.trashinfo")
	if !success {
		t.Error("file has not been deleted from trash")
	}
}

func TestRestoreFromTrash(t *testing.T) {}

func TestEmptyTrash(t *testing.T) {
	trashPath := "/home/user/.local/share/Trash"
	createTrashFile("script.sh")
	createTrashFile("lib.so")
	EmptyTrash()
	success := !existsTrashFile(trashPath, "script.sh") && !existsTrashFile(trashPath, "lib.so")
	if !success {
		t.Error("trash has not been emptied")
	}
}

func createTrashFile(filename string) {
	trashPath := "/home/user/.local/share/Trash"
	fs.MkdirAll(trashPath+"/files", os.ModeDir)
	fs.Create(trashPath + "/files/" + filename)
	fs.MkdirAll(trashPath+"/info", os.ModeDir)
	fs.Create(trashPath + "/info/" + filename + ".trashinfo")
}

func existsTrashFile(trashPath string, filename string) bool {
	return exists(trashPath+"/files/"+filename) && exists(trashPath+"/info/"+filename+".trashinfo")
}

func initHomeEnvironment() {
	fs.MkdirAll("/home/user/.local/share/Trash", os.ModeDir)
	fs.Create("/home/user/file")
}

func exists(path string) bool {
	dir, err := fs.Stat(path)
	return err == nil && dir.Mode().IsDir()
}
