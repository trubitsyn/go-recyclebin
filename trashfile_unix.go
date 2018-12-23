// +build !windows

package recyclebin

import (
	"github.com/spf13/afero"
	"path"
	"strconv"
	"strings"
)

func getTrashedFilename(trashPath string, filename string) string {
	isDuplicateFilename, _ := afero.Exists(fs, buildTrashFilePath(trashPath, filename))
	if isDuplicateFilename {
		filename = generateNewFilename(filename)
	}
	return filename
}

func generateNewFilename(existingFilename string) string {
	extension := path.Ext(existingFilename)
	bareName := strings.TrimSuffix(existingFilename, extension)
	newFilename := existingFilename
	index := -1

	isDuplicateFilename, _ := afero.Exists(fs, newFilename)
	for index == -1 || isDuplicateFilename {
		index += 1
		newFilename = bareName + strconv.Itoa(index) + extension
	}
	return newFilename
}

func buildTrashFilePath(trashPath string, filename string) string {
	return trashPath + "/files/" + filename
}
