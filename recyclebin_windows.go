// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to remove zero-width characters from strings.
package recyclebin

// #include "recyclebin.h"
import "C"

func MoveToTrash(filepath string) error {
	C.move_to_trash(filepath)
}

func RestoreFromTrash(filename string) error {
	C.restore_from_trash(filename)
	return nil
}

func DeleteFromTrash(filename string) error {
	C.delete_from_trash(filename)
	return nil
}

func EmptyTrash() error {
	C.empty_trash()
	return nil
}
