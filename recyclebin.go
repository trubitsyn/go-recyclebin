// Copyright 2018 Nikola Trubitsyn. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recyclebin implements functions to use Trash (or Recycle Bin).
package recyclebin

import (
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

type RecycleBin interface {
	Recycle(filename string) error
	Restore(trashFilename string) error
	Remove(trashFilename string) error
	Empty() error
}
