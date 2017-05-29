package archive

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// For testing; we can replace these with versions that intercept calls as we need.
var create = os.Create
var open = os.Open
var readDir = ioutil.ReadDir
var stat = os.Stat

// Dirsf takes an array of directory paths as strings, a formatting string for the file names, and a
// FormatName to identify the type of archive to produce; and produces archives for each of the
// given directories. If any of the directory names don't exist or aren't directories, an error will
// be returned.
//
// The values in `dirnames` can be absolute, or relative paths for the directories. These are simply
// passed into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%d` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%d"`.
//
// Upon success, an array of the archive filenames will be returned.
func Dirsf(dirnames []string, nameFmt string, formatName FormatName) ([]string, error) {
	// Cores is the number of logical CPU cores the Go runtime has available to it.
	cores := runtime.GOMAXPROCS(0)

	limChan := make(chan bool, cores)
	errChan := make(chan error, len(dirnames))
	resChan := make(chan string, len(dirnames))

	// Prepare the limiter. We fill the channel with as many values as we want archives to be
	// created concurrently; for now, this is the number of logical CPU cores available.
	for i := 0; i < cores; i++ {
		limChan <- true
	}

	for i, dirname := range dirnames {
		<-limChan

		go func(i int, dirname string) {
			log.Printf("Started archiving directory '%s'...", dirname)

			res, err := Dirf(dirname, nameFmt, formatName)
			if err != nil {
				errChan <- err
			} else {
				resChan <- res
			}

			log.Printf("Finished archiving directory '%s'...", dirname)

			// Release use of limiter
			limChan <- true
		}(i, dirname)
	}

	filenames := []string{}

	for {
		select {
		case err := <-errChan:
			return []string{}, err
		case res := <-resChan:
			filenames = append(filenames, res)
		}

		if len(dirnames) == len(filenames) {
			break
		}
	}

	sort.Strings(filenames)

	return filenames, nil
}

// Dirf archives a given source directory, and creates an archive with a name in the given format,
// in the given archive artifact format (FormatName). If the dirname given does not exist, or is not
// a directory, an error will be returned.
//
// The value of `dirname` can be an absolute or relative path to a directory. It is simply passed
// into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%c` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%c"`.
//
// Upon success, the archive filename will be returned.
func Dirf(dirname string, nameFmt string, formatName FormatName) (string, error) {
	parentPath := path.Dir(dirname)

	// Create the destination filename based on the name format, and base path.
	fileName := fmt.Sprintf(nameFmt, path.Base(dirname), time.Now().Unix())
	fileName = strings.Replace(fileName, " ", "_", -1)

	format, err := findFormatByName(formatName)
	if err != nil {
		return "", err
	}

	// Produce the archive file, with the given name, in the given directory.
	artifact, err := format.producer(parentPath, fileName)
	if err != nil {
		return "", err
	}

	defer artifact.Close()

	return artifact.Name(), walk(dirname, artifact)
}

func walk(root string, artifact Artifact) error {
	info, err := stat(root)
	if err != nil {
		return err
	}

	return doWalk(root, info, artifact)
}

// walk traverses a directory tree, starting at the given path. This is a simplified version of the
// walk function provided in the standard library designed to make testing a little easier.
func doWalk(path string, info os.FileInfo, artifact Artifact) error {
	err := artifact.AddFile(path, info)
	if err != nil {
		return err
	}

	// Bail if we're not looking at a directory, we have nothing left to do.
	if !info.IsDir() {
		return nil
	}

	// Read all of the files in this directory.
	files, err := readDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Create the full path to file.
		filename := filepath.Join(path, file.Name())

		err = doWalk(filename, file, artifact)
		if err != nil {
			return err
		}
	}

	return nil
}
