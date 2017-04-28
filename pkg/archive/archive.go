package archive

import (
	"fmt"
	"path"
	"runtime"
	"sort"
	"time"
)

// cores is the number of logical CPU cores the Go runtime has available to it.
var cores = runtime.GOMAXPROCS(0)

// Dirsf takes an array of directory paths as strings, and a formatting string for the file
// names, and produces .tar.gz archives for each of the given directories. If any of the directories
// don't exist, an error will be returned.
//
// The values in `dirnames` can be absolute, or relative paths for the directories. These are simply
// passed into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%d` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%d"`.
//
// Upon success, an array of the archive filenames will be returned.
func Dirsf(dirnames []string, namefmt string) ([]string, error) {
	tarFiles := []string{}
	timestamp := time.Now().Unix()

	// Create the names for each archive that will be produced, based on the format provided.
	namefmt = fmt.Sprintf("%s.tar.gz", namefmt)

	limiter := make(chan bool, cores)

	// Set up the limiter, we fill it with the number of cores we have available.
	for i := 0; i < cores; i++ {
		limiter <- true
	}

	for _, dirname := range dirnames {
		<-limiter

		go func(d string) {
			// Simulate longer process.
			time.Sleep(2 * time.Second)
			fmt.Printf("%s (%d)\n", d, time.Now().Unix())

			// Release "worker" availability
			limiter <- true
		}(dirname)
	}

	// Wait for all workers to finish.
	for i := 0; i < cores; i++ {
		<-limiter
	}

	// ...

	// We probably want to batch this process, we should have that code lying around somewhere still
	// with channels and threads and whatnot...
	for _, dirname := range dirnames {
		basename := path.Base(dirname)

		tarFiles = append(tarFiles, fmt.Sprintf(namefmt, basename, timestamp))
	}

	fmt.Println(cores)

	// Sort tar filenames so they still come out in alphabetical order.
	sort.Strings(tarFiles)

	return tarFiles, nil
}

// @todo: Instead of splitting the data into chunks, can we process one dir at a time up to a limit
// (the limit being the number of cores available to the Go runtime).
//   This also has the benefit that if one core is faster than others for some reason, it won't sit
//   idle while others have still handling their "chunk".
//
//   Can you count the length of a channel? If so, we should add to the channel of ongoing tasks
//   before we start the Go routine for it, otherwise we could have race conditions.
//
// @todo: Can this be a pipeline'd process using channels?
//   Maybe not, given it's going to be doing IO, we can use channels to inform the application about
//   when something finishes though.
//
// @todo: How on earth do we test this?
//   Internal implementation maybe doesn't matter? As long we're blocking waiting for all of them to
//   finish then the tests should check the result, not what the process is actually doing.
//
// @todo: Do we care about outputting anything during this process? Reporting on progress?
//   Maybe some logging that is spread application-wide is needed? Or maybe just in the command to
//   be honest, it's not like you should be sat watching foldup do it's thing.
//
// @todo: How do we handle failure with in the Go routines? Can we cancel all of the processes? And
// if we do this, do we need to clean up?
//   ???
//
