// +build !windows

package recyclebin

import (
	"github.com/spf13/afero"
	"path"
	"strconv"
	"strings"
)

func getTrashedFilename(trashPath string, filename string) string {
	extension := path.Ext(filename)
	bareName := strings.TrimSuffix(filename, extension)
	newFilename := filename
	index := -1
	isDuplicateFilename := true
	for isDuplicateFilename {
		index += 1
		newFilename = bareName + strconv.Itoa(index) + extension
		existsTrashFile, _ := afero.Exists(fs, buildTrashFilePath(trashPath, newFilename))
		existsTrashInfo, _ := afero.Exists(fs, buildTrashInfoPath(trashPath, newFilename))
		isDuplicateFilename = existsTrashFile || existsTrashInfo
	}
	return newFilename
}

func buildTrashFilePath(trashPath string, filename string) string {
	return trashPath + "/files/" + filename
}
