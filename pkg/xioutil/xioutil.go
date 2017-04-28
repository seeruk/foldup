package xioutil

import (
	"io/ioutil"
	"os"
)

// ReadDirsInDir reads the directory named by dirname and returns a list of directory entries sorted
// by filename.
func ReadDirsInDir(dirname string, hidden bool) ([]os.FileInfo, error) {
	dirs := []os.FileInfo{}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return dirs, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		name := f.Name()

		// If we're not showing hidden folders, let's skip those too.
		if !hidden && len(name) > 0 && name[:1] == "." {
			continue
		}

		dirs = append(dirs, f)
	}

	return dirs, nil
}
