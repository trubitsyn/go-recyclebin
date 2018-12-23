// +build !windows

package recyclebin

import (
	"github.com/spf13/afero"
	"path"
	"strconv"
	"strings"
)

func getTrashedFilename(trashPath string, filename string) string {
	trashedFilename := filename
	isDuplicateFilename, _ := afero.Exists(fs, trashPath+"/files/"+trashedFilename)
	if isDuplicateFilename {
		trashedFilename = generateNewFilename(trashedFilename)
	}
	return trashedFilename
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

func buildTrashFilePath(trashInfoFilePath string) (string, error) {
	trashInfo, err := readTrashInfo(trashInfoFilePath)
	if err != nil {
		return "", err
	}
	return trashInfo.Path, nil
}
