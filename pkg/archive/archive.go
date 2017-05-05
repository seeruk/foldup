package archive

import (
	"fmt"
	"io/ioutil"
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

// Dirsf takes an array of directory paths as strings, and a formatting string for the file names,
// and produces archives for each of the given directories. If any of the directory names don't
// exist or aren't directories, an error will be returned.
//
// The values in `dirnames` can be absolute, or relative paths for the directories. These are simply
// passed into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%c` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%c"`.
//
// Upon success, an array of the archive filenames will be returned.
func Dirsf(dirnames []string, nameFmt string, format Format) ([]string, error) {
	// Cores is the number of logical CPU cores the Go runtime has available to it.
	cores := runtime.GOMAXPROCS(0)

	limiter := make(chan bool, cores)
	errs := make(chan error, 1)

	// Prepare the limiter. We fill the channel with as many values as we want archives to be
	// created concurrently; for now, this is the number of logical CPU cores available.
	for i := 0; i < cores; i++ {
		limiter <- true
	}

	filenames := make([]string, len(dirnames))

	// artifact each directory, if any one fails, we stop then and return that first error.
	for i, dirname := range dirnames {
		select {
		case <-limiter:
			// Process archiving a directory asynchronously.
			go func(i int, dirname string) {
				filename, err := Dirf(dirname, nameFmt, format)
				if err != nil {
					errs <- err
				}

				// @todo: Should we only do this if we didn't have an error?
				filenames[i] = filename

				// Release use of limiter
				limiter <- true
			}(i, dirname)
		case err := <-errs:
			// @todo: Should we just log errors and carry on? This will halt the entire backup
			// currently. Either way it's not good I guess...
			//
			// If we do carry on, then it will probably mean that we can clean up better though
			// later on if we continue to try upload, or if we just want to delete the other
			// produced archives so we're not taking up a bunch of disk space.
			if err != nil {
				return filenames, err
			}
		}
	}

	// Wait for all workers to finish, if this blocks then something has gone quite wrong.
	for i := 0; i < cores; i++ {
		<-limiter
	}

	sort.Strings(filenames)

	return filenames, nil
}

// Dirf archives a given source directory, and creates an archive with a name in the given format.
// If the dirname given does not exist, or is not a directory, an error will be returned.
//
// The value of `dirname` can be an absolute or relative path to a directory. It is simply passed
// into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%c` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%c"`.
//
// Upon success, the archive filename will be returned.
func Dirf(dirname string, nameFmt string, format Format) (string, error) {
	parentPath := path.Dir(dirname)

	// Create the destination filename based on the name format, and base path.
	fileName := fmt.Sprintf(nameFmt, path.Base(dirname), time.Now().Unix())
	fileName = strings.Replace(fileName, " ", "_", -1)

	producer := producers[format]

	// Produce the archive file, with the given name, in the given directory.
	artifact, err := producer(parentPath, fileName)
	if err != nil {
		return "", err
	}

	defer artifact.Close()

	// @todo: Could we put an unexposed walk function in this package? That would leave this file
	// with only a single remaining usage of xos.FileSystem - which could also be eliminated.
	// @todo: How do we get the name from the archive?
	return artifact.Name(), walk(dirname, artifact)
}

func walk(root string, artifact artifact) error {
	info, err := stat(root)
	if err != nil {
		return err
	}

	return doWalk(root, info, artifact)
}

// walk traverses a directory tree, starting at the given path. This is a simplified version of the
// walk function provided in the standard library designed to make testing a little easier.
func doWalk(path string, info os.FileInfo, artifact artifact) error {
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
