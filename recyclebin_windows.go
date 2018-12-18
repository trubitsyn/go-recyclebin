// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to remove zero-width characters from strings.
package recyclebin

// #include "recyclebin.h"
import "C"

type windowsRecycleBin struct {
}

func (bin windowsRecycleBin) Recycle(filepath string) {
	C.move_to_trash(filepath)
}

func (bin windowsRecycleBin) Restore(trashFilename string) error {
	C.restore_from_trash(trashFilename)
	return nil
}

func (bin windowsRecycleBin) Delete(trashFilename string) error {
	C.delete_from_trash(trashFilename)
	return nil
}

func (bin windowsRecycleBin) Empty() error {
	C.empty_trash()
	return nil
}
